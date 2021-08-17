// *** WARNING: this file was generated by tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package codegentest

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The response from the ListKeys operation.
// API Version: 2021-02-01.
func ListStorageAccountKeys(ctx *pulumi.Context, args *ListStorageAccountKeysArgs, opts ...pulumi.InvokeOption) (*ListStorageAccountKeysResult, error) {
	var rv ListStorageAccountKeysResult
	err := ctx.Invoke("azure-native:codegentest:listStorageAccountKeys", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

type ListStorageAccountKeysArgs struct {
	// The name of the storage account within the specified resource group. Storage account names must be between 3 and 24 characters in length and use numbers and lower-case letters only.
	AccountName string `pulumi:"accountName"`
	// Specifies type of the key to be listed. Possible value is kerb.
	Expand *string `pulumi:"expand"`
	// The name of the resource group within the user's subscription. The name is case insensitive.
	ResourceGroupName string `pulumi:"resourceGroupName"`
}


// The response from the ListKeys operation.
type ListStorageAccountKeysResult struct {
	// Gets the list of storage account keys and their properties for the specified storage account.
	Keys []StorageAccountKeyResponse `pulumi:"keys"`
}


func ListStorageAccountKeysOutput(ctx *pulumi.Context, args ListStorageAccountKeysOutputArgs, opts ...pulumi.InvokeOption) ListStorageAccountKeysResultTypeOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (ListStorageAccountKeysResult, error) {
			args := v.(ListStorageAccountKeysArgs)
			r, err := ListStorageAccountKeys(ctx, &args, opts...)
			return *r, err
		}).(ListStorageAccountKeysResultTypeOutput)
}

type ListStorageAccountKeysOutputArgs struct {
	// The name of the storage account within the specified resource group. Storage account names must be between 3 and 24 characters in length and use numbers and lower-case letters only.
	AccountName pulumi.StringInput `pulumi:"accountName"`
	// Specifies type of the key to be listed. Possible value is kerb.
	Expand pulumi.StringPtrInput `pulumi:"expand"`
	// The name of the resource group within the user's subscription. The name is case insensitive.
	ResourceGroupName pulumi.StringInput `pulumi:"resourceGroupName"`
}

func (ListStorageAccountKeysOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*ListStorageAccountKeysArgs)(nil)).Elem()
}

// The response from the ListKeys operation.
type ListStorageAccountKeysResultTypeOutput struct { *pulumi.OutputState }

func (ListStorageAccountKeysResultTypeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*ListStorageAccountKeysResult)(nil)).Elem()
}

func (o ListStorageAccountKeysResultTypeOutput) ToListStorageAccountKeysResultTypeOutput() ListStorageAccountKeysResultTypeOutput {
	return o
}

func (o ListStorageAccountKeysResultTypeOutput) ToListStorageAccountKeysResultTypeOutputWithContext(ctx context.Context) ListStorageAccountKeysResultTypeOutput {
	return o
}

// Gets the list of storage account keys and their properties for the specified storage account.
func (o ListStorageAccountKeysResultTypeOutput) Keys() StorageAccountKeyResponseArrayOutput {
	return o.ApplyT(func (v ListStorageAccountKeysResult) []StorageAccountKeyResponse { return v.Keys }).(StorageAccountKeyResponseArrayOutput)
}


func init() {
        pulumi.RegisterOutputType(ListStorageAccountKeysResultTypeOutput{})
}

