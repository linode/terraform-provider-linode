import boto3
import sys
import os
import xml.etree.ElementTree as ET
from botocore.exceptions import NoCredentialsError

ACCESS_KEY = os.environ.get('LINODE_CLI_OBJ_ACCESS_KEY')
SECRET_KEY = os.environ.get('LINODE_CLI_OBJ_SECRET_KEY')
BUCKET_NAME = 'dx-test-results'

linode_obj_config = {
    "aws_access_key_id": ACCESS_KEY,
    "aws_secret_access_key": SECRET_KEY,
    "endpoint_url": "https://us-southeast-1.linodeobjects.com",
    "region_name": "us-southeast-1",
}

def change_xml_report_to_tod_acceptable_version(file_name):
    # Load the original XML file
    tree = ET.parse(file_name)
    root = tree.getroot()

    testsuites_element = root

    # total
    total_tests = int(testsuites_element.get('tests'))
    total_failures = int(testsuites_element.get('failures'))
    total_errors = int(testsuites_element.get('errors'))
    total_skipped = int(testsuites_element.get('skipped'))

    # Create a new <testsuites> element with aggregated values
    new_testsuites = ET.Element("testsuites")
    new_testsuites.set("tests", str(total_tests))
    new_testsuites.set("failures", str(total_failures))
    new_testsuites.set("errors", str(total_errors))
    new_testsuites.set("skipped", str(total_skipped))

    # Create a new <testsuite> element under <testsuites>
    new_testsuite = ET.SubElement(new_testsuites, "testsuite", attrib=testsuites_element.attrib)

    for testcase in root.findall('.//testcase'):
        new_testcase = ET.SubElement(new_testsuite, "testcase", attrib=testcase.attrib)
        for child in testcase:
            new_testcase.append(child)

    # Save the new XML to a file
    new_tree = ET.ElementTree(new_testsuites)
    try:
        new_tree = ET.ElementTree(root)
        new_tree.write(file_name, encoding="UTF-8", xml_declaration=True)
        print("XML content successfully over-written to " + file_name)

    except Exception as e:
        print("Error writing XML content:", str(e))



def upload_to_linode_object_storage(file_name):
    try:
        s3 = boto3.client('s3', **linode_obj_config)

        s3.upload_file(Filename=file_name, Bucket=BUCKET_NAME, Key=file_name)

        print(f'Successfully uploaded {file_name} to Linode Object Storage.')

    except NoCredentialsError:
        print('Credentials not available. Ensure you have set your AWS credentials.')


if __name__ == '__main__':
    if len(sys.argv) != 2:
        print('Usage: python upload_to_linode.py <file_name>')
        sys.exit(1)

    file_name = sys.argv[1]

    if not file_name:
        print('Error: The provided file name is empty or invalid.')
        sys.exit(1)

    change_xml_report_to_tod_acceptable_version(file_name)
    upload_to_linode_object_storage(file_name)