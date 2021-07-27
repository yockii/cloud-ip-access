package service

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	rds20140815 "github.com/alibabacloud-go/rds-20140815/v2/client"
	teaservice "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/yockii/cloud-ip-access/util"
)

var RdsGroupService = new(rdsGroupService)

type rdsGroupService struct {
	aliyunClient *rds20140815.Client
}

func init() {
	config := &openapi.Config{
		AccessKeyId:     tea.String(util.Config.GetString("aliyun.accessKey")),
		AccessKeySecret: tea.String(util.Config.GetString("aliyun.accessSecret")),
	}
	// 访问的域名
	config.Endpoint = tea.String("rds.aliyuncs.com")
	RdsGroupService.aliyunClient, _ = rds20140815.NewClient(config)
}

func (r *rdsGroupService) UpdateRdsDBInstancesWhiteIps(ip string) error {
	dbList, err := r.FindRdsDBInstances()
	if err != nil {
		return err
	}
	for _, db := range dbList {
		ipArray, err := r.FindRdsDBInstanceIpWhiteList(db.DBInstanceId)
		if err != nil {
			return err
		}
		var modifyIp *rds20140815.DescribeDBInstanceIPArrayListResponseBodyItemsDBInstanceIPArray
		for _, dbIp := range ipArray {
			if *dbIp.DBInstanceIPArrayName == "company" {
				modifyIp = dbIp
				break
			}
		}
		if modifyIp != nil {
			if err := r.UpdateRdsDbInstanceIpWhiteList(ip, db.DBInstanceId, modifyIp.DBInstanceIPArrayName); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *rdsGroupService) FindRdsDBInstances() ([]*rds20140815.DescribeDBInstancesResponseBodyItemsDBInstance, error) {
	describeDBInstancesRequest := &rds20140815.DescribeDBInstancesRequest{
		RegionId: tea.String("cn-hangzhou"),
	}
	resp, err := r.aliyunClient.DescribeDBInstancesWithOptions(describeDBInstancesRequest,
		&teaservice.RuntimeOptions{
			IgnoreSSL: tea.Bool(true),
		})
	if err != nil {
		return nil, err
	}
	return resp.Body.Items.DBInstance, nil
}

func (r *rdsGroupService) FindRdsDBInstanceIpWhiteList(instanceId *string) ([]*rds20140815.DescribeDBInstanceIPArrayListResponseBodyItemsDBInstanceIPArray, error) {
	describeDBInstanceIPArrayListRequest := &rds20140815.DescribeDBInstanceIPArrayListRequest{
		DBInstanceId: instanceId,
	}
	resp, err := r.aliyunClient.DescribeDBInstanceIPArrayListWithOptions(describeDBInstanceIPArrayListRequest,
		&teaservice.RuntimeOptions{
			IgnoreSSL: tea.Bool(true),
		})
	if err != nil {
		return nil, err
	}
	return resp.Body.Items.DBInstanceIPArray, nil
}

func (r *rdsGroupService) UpdateRdsDbInstanceIpWhiteList(ip string, instanceId, ipArrayName *string) error {
	modifySecurityIpsRequest := &rds20140815.ModifySecurityIpsRequest{
		DBInstanceId:          instanceId,
		SecurityIps:           tea.String(ip),
		DBInstanceIPArrayName: ipArrayName,
	}
	_, err := r.aliyunClient.ModifySecurityIpsWithOptions(modifySecurityIpsRequest,
		&teaservice.RuntimeOptions{
			IgnoreSSL: tea.Bool(true),
		})
	return err
}
