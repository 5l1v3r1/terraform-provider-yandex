package yandex

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/genproto/protobuf/field_mask"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/iam/v1"
)

const yandexIAMServiceAccountDefaultTimeout = 1 * time.Minute

func resourceYandexIAMServiceAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceYandexIAMServiceAccountCreate,
		Read:   resourceYandexIAMServiceAccountRead,
		Update: resourceYandexIAMServiceAccountUpdate,
		Delete: resourceYandexIAMServiceAccountDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(yandexIAMServiceAccountDefaultTimeout),
			Update: schema.DefaultTimeout(yandexIAMServiceAccountDefaultTimeout),
			Delete: schema.DefaultTimeout(yandexIAMServiceAccountDefaultTimeout),
		},

		SchemaVersion: 0,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"folder_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceYandexIAMServiceAccountCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	folderID, err := getFolderID(d, config)
	if err != nil {
		return fmt.Errorf("Error getting folder ID while creating service account: %s", err)
	}

	req := iam.CreateServiceAccountRequest{
		FolderId:    folderID,
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()

	op, err := config.sdk.WrapOperation(config.sdk.IAM().ServiceAccount().Create(ctx, &req))
	if err != nil {
		return fmt.Errorf("Error while requesting API to create service account: %s", err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error while waiting operation to create service account: %s", err)
	}

	resp, err := op.Response()
	if err != nil {
		return fmt.Errorf("Service account creation failed: %s", err)
	}

	sa, ok := resp.(*iam.ServiceAccount)
	if !ok {
		return fmt.Errorf("Create response doesn't contain Service Account")
	}

	d.SetId(sa.Id)

	return resourceYandexIAMServiceAccountRead(d, meta)
}

func resourceYandexIAMServiceAccountRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sa, err := config.sdk.IAM().ServiceAccount().Get(context.Background(), &iam.GetServiceAccountRequest{
		ServiceAccountId: d.Id(),
	})

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account %q", d.Get("name").(string)))
	}

	createdAt, err := getTimestamp(sa.CreatedAt)
	if err != nil {
		return err
	}

	d.Set("created_at", createdAt)
	d.Set("name", sa.Name)
	d.Set("folder_id", sa.FolderId)
	d.Set("description", sa.Description)

	return nil
}

func resourceYandexIAMServiceAccountUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	req := &iam.UpdateServiceAccountRequest{
		ServiceAccountId: d.Id(),
		UpdateMask:       &field_mask.FieldMask{},
	}

	if d.HasChange("name") {
		req.Name = d.Get("name").(string)
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "name")
	}

	if d.HasChange("description") {
		req.Description = d.Get("description").(string)
		req.UpdateMask.Paths = append(req.UpdateMask.Paths, "description")
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()

	op, err := config.sdk.WrapOperation(config.sdk.IAM().ServiceAccount().Update(ctx, req))
	if err != nil {
		return fmt.Errorf("Error while requesting API to update Service Account %q: %s", d.Id(), err)
	}

	err = op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error updating Service Account %q: %s", d.Id(), err)
	}

	for _, v := range req.UpdateMask.Paths {
		d.SetPartial(v)
	}

	d.Partial(false)

	return resourceYandexIAMServiceAccountRead(d, meta)
}

func resourceYandexIAMServiceAccountDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[DEBUG] Deleting Service Account %q", d.Id())

	req := &iam.DeleteServiceAccountRequest{
		ServiceAccountId: d.Id(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()

	op, err := config.sdk.WrapOperation(config.sdk.IAM().ServiceAccount().Delete(ctx, req))
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Service Account %q", d.Get("name").(string)))
	}

	err = op.Wait(ctx)
	if err != nil {
		return err
	}

	resp, err := op.Response()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Finished deleting Service Account %q: %#v", d.Id(), resp)
	return nil
}
