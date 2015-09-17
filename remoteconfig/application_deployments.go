package remoteconfig

func (s *ApplicationDeployments) ListApps() (appNameList []string) {
	for k, _ := range s.Applications {
		appNameList = append(appNameList, k)
	}
	return
}

func (s *ApplicationDeployments) ListDeployments(appName string) (deploymentList []string) {
	for _, k := range s.Applications[appName].Deployments {
		deploymentList = append(deploymentList, k.Name)
	}
	return
}

func (s *ApplicationDeployments) GetDeployment(appName, deploymentName string) (deployment Deployment) {
	for _, d := range s.Applications[appName].Deployments {

		if d.Name == deploymentName {
			deployment = d
		}
	}
	return
}
