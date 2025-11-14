//go:build integration || producerimagesharegroup

package producerimagesharegroup_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/linode/terraform-provider-linode/v3/linode/acceptance"
	"github.com/linode/terraform-provider-linode/v3/linode/producerimagesharegroup/tmpl"
)

func TestAccResourceImageShareGroup_basic(t *testing.T) {
	t.Parallel()

	resourceName := "linode_producer_image_share_group.foobar"
	label := acctest.RandomWithPrefix("tf-test")
	description := "A cool description."

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, label, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "label", label),
					resource.TestCheckResourceAttr(resourceName, "description", description),
					resource.TestCheckResourceAttr(resourceName, "is_suspended", "false"),
					resource.TestCheckResourceAttr(resourceName, "images_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "members_count", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckNoResourceAttr(resourceName, "updated"),
					resource.TestCheckNoResourceAttr(resourceName, "expiry"),
				),
			},
		},
	})
}

func TestAccResourceImageShareGroup_updates(t *testing.T) {
	t.Parallel()

	resourceName := "linode_producer_image_share_group.foobar"
	label := acctest.RandomWithPrefix("tf-test")
	testRegion, err := acceptance.GetRandomRegionWithCaps([]string{"Linodes"}, "core")
	if err != nil {
		t.Fatalf("failed to get test region: %s", err)
	}
	imageLabel1 := "image" + acctest.RandomWithPrefix("tf-test")
	imageLabel2 := "image" + acctest.RandomWithPrefix("tf-test")
	isgLabel := "sharegroup" + acctest.RandomWithPrefix("tf-test")
	isgDescription := "A cool description."
	isgLabelUpdated := isgLabel + "-updated"
	isgDescriptionUpdated := isgDescription + " updated"

	imagesStep2 := []tmpl.ShareGroupImageTemplate{
		{
			ID:          "${linode_image.foobar.id}",
			Label:       "Share-Image-1",
			Description: "Share Image 1 Description",
		},
	}

	imagesStep3 := []tmpl.ShareGroupImageTemplate{
		{
			ID:          "${linode_image.foobar.id}",
			Label:       "Share-Image-1-updated",
			Description: "Share Image 1 Description updated",
		},
		{
			ID:          "${linode_image.barfoo.id}",
			Label:       "Share-Image-2",
			Description: "Share Image 2 Description",
		},
	}

	imagesStep4 := []tmpl.ShareGroupImageTemplate{
		{
			ID:          "${linode_image.foobar.id}",
			Label:       "Share-Image-1-updated",
			Description: "Share Image 1 Description updated",
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.PreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create empty share group
			{
				Config: tmpl.Updates(t, label, testRegion, imageLabel1, imageLabel2, isgLabel, isgDescription, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "label", isgLabel),
					resource.TestCheckResourceAttr(resourceName, "description", isgDescription),
					resource.TestCheckResourceAttr(resourceName, "is_suspended", "false"),
					resource.TestCheckResourceAttr(resourceName, "images_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "members_count", "0"),
				),
			},
			// Step 2: Add first image
			{
				Config: tmpl.Updates(t, label, testRegion, imageLabel1, imageLabel2, isgLabel, isgDescription, imagesStep2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "label", isgLabel),
					resource.TestCheckResourceAttr(resourceName, "description", isgDescription),
					resource.TestCheckResourceAttr(resourceName, "is_suspended", "false"),
					resource.TestCheckResourceAttr(resourceName, "images_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "members_count", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.id"),
					resource.TestCheckResourceAttr(resourceName, "images.0.label", "Share-Image-1"),
					resource.TestCheckResourceAttr(resourceName, "images.0.description", "Share Image 1 Description"),
				),
			},
			// Step 3: Add second image and update first image
			{
				Config: tmpl.Updates(t, label, testRegion, imageLabel1, imageLabel2, isgLabel, isgDescription, imagesStep3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "label", isgLabel),
					resource.TestCheckResourceAttr(resourceName, "description", isgDescription),
					resource.TestCheckResourceAttr(resourceName, "is_suspended", "false"),
					resource.TestCheckResourceAttr(resourceName, "images_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "members_count", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.id"),
					resource.TestCheckResourceAttr(resourceName, "images.0.label", "Share-Image-1-updated"),
					resource.TestCheckResourceAttr(resourceName, "images.0.description", "Share Image 1 Description updated"),
					resource.TestCheckResourceAttrSet(resourceName, "images.1.id"),
					resource.TestCheckResourceAttr(resourceName, "images.1.label", "Share-Image-2"),
					resource.TestCheckResourceAttr(resourceName, "images.1.description", "Share Image 2 Description"),
				),
			},
			// Step 4: Update the Share Group and remove the second image
			{
				Config: tmpl.Updates(t, label, testRegion, imageLabel1, imageLabel2, isgLabelUpdated, isgDescriptionUpdated, imagesStep4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "label", isgLabel+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", isgDescription+" updated"),
					resource.TestCheckResourceAttr(resourceName, "is_suspended", "false"),
					resource.TestCheckResourceAttr(resourceName, "images_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "members_count", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "images.0.id"),
					resource.TestCheckResourceAttr(resourceName, "images.0.label", "Share-Image-1-updated"),
					resource.TestCheckResourceAttr(resourceName, "images.0.description", "Share Image 1 Description updated"),
				),
			},
			// Step 5: Remove the first image
			{
				Config: tmpl.Updates(t, label, testRegion, imageLabel1, imageLabel2, isgLabelUpdated, isgDescriptionUpdated, nil),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "label", isgLabel+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", isgDescription+" updated"),
					resource.TestCheckResourceAttr(resourceName, "is_suspended", "false"),
					resource.TestCheckResourceAttr(resourceName, "images_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "members_count", "0"),
				),
			},
		},
	})
}
