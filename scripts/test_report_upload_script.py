import boto3
import sys
import os
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

    upload_to_linode_object_storage(file_name)