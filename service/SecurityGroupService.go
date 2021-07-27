package service

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	ecs20140526 "github.com/alibabacloud-go/ecs-20140526/v2/client"
	teaservice "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/yockii/cloud-ip-access/util"
)

var SecurityGroupService = new(securityGroupService)

type securityGroupService struct {
	aliyunClient *ecs20140526.Client
}

func init() {
	config := &openapi.Config{
		AccessKeyId:     tea.String(util.Config.GetString("aliyun.accessKey")),
		AccessKeySecret: tea.String(util.Config.GetString("aliyun.accessSecret")),
	}
	// 访问的域名
	config.Endpoint = tea.String("ecs-cn-hangzhou.aliyuncs.com")
	SecurityGroupService.aliyunClient, _ = ecs20140526.NewClient(config)
}

func (s *securityGroupService) UpdateSecurityGroupManageIP(ip string) error {
	sgList, err := s.FindSecurityGroups()
	if err != nil {
		return err
	}
	var updatedSg *ecs20140526.DescribeSecurityGroupsResponseBodySecurityGroupsSecurityGroup
	for _, sg := range sgList {
		if *sg.SecurityGroupName == "manage" {
			updatedSg = sg
			break
		}
	}
	if updatedSg != nil {
		// 查找安全组的规则
		pList, err := s.FindSecurityGroupRules(updatedSg.SecurityGroupId)
		if err != nil {
			return err
		}
		for _, p := range pList {
			if *p.Description == "公司IP" {
				if err2 := s.ModifySecurityGroup(ip, updatedSg.SecurityGroupId, p); err2 != nil {
					return err2
				}
				break
			}
		}
	}
	return nil
}

func (s *securityGroupService) FindSecurityGroups() ([]*ecs20140526.DescribeSecurityGroupsResponseBodySecurityGroupsSecurityGroup, error) {
	resp, err := s.aliyunClient.DescribeSecurityGroupsWithOptions(&ecs20140526.DescribeSecurityGroupsRequest{
		RegionId: tea.String("cn-hangzhou"),
	},
		&teaservice.RuntimeOptions{
			IgnoreSSL: tea.Bool(true),
		},
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return resp.Body.SecurityGroups.SecurityGroup, nil
}

func (s *securityGroupService) FindSecurityGroupRules(sgId *string) ([]*ecs20140526.DescribeSecurityGroupAttributeResponseBodyPermissionsPermission, error) {
	describeSecurityGroupAttributeRequest := &ecs20140526.DescribeSecurityGroupAttributeRequest{
		SecurityGroupId: sgId,
		RegionId:        tea.String("cn-hangzhou"),
	}
	resp, err := s.aliyunClient.DescribeSecurityGroupAttributeWithOptions(describeSecurityGroupAttributeRequest,
		&teaservice.RuntimeOptions{
			IgnoreSSL: tea.Bool(true),
		})
	if err != nil {
		return nil, err
	}
	return resp.Body.Permissions.Permission, nil
}

func (s *securityGroupService) ModifySecurityGroup(ip string, sgId *string, p *ecs20140526.DescribeSecurityGroupAttributeResponseBodyPermissionsPermission) error {
	revokeSecurityGroupRequest := &ecs20140526.RevokeSecurityGroupRequest{
		RegionId:        tea.String("cn-hangzhou"),
		SecurityGroupId: sgId,
		PortRange:       p.PortRange,
		IpProtocol:      p.IpProtocol,
		SourceCidrIp:    p.SourceCidrIp,
	}
	_, err := s.aliyunClient.RevokeSecurityGroupWithOptions(revokeSecurityGroupRequest,
		&teaservice.RuntimeOptions{
			IgnoreSSL: tea.Bool(true),
		})
	if err != nil {
		return err
	}

	authorizeSecurityGroupRequest := &ecs20140526.AuthorizeSecurityGroupRequest{
		RegionId:        tea.String("cn-hangzhou"),
		SecurityGroupId: sgId,
		IpProtocol:      p.IpProtocol,
		PortRange:       p.PortRange,
		SourceCidrIp:    tea.String(ip),
		Description:     p.Description,
	}
	_, err = s.aliyunClient.AuthorizeSecurityGroupWithOptions(authorizeSecurityGroupRequest,
		&teaservice.RuntimeOptions{
			IgnoreSSL: tea.Bool(true),
		})
	return err
}
