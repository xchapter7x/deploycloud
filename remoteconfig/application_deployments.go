package remoteconfig

//ListApps - list the apps available in the configuration file
func (s *ApplicationDeployments) ListApps() (appNameList []string) {
	for k, _ := range s.Applications {
		appNameList = append(appNameList, k)
	}
	return
}

//ListDeployments - list the deployments available for the given application
func (s *ApplicationDeployments) ListDeployments(appName string) (deploymentList []string) {
	for _, k := range s.Applications[appName].Deployments {
		deploymentList = append(deploymentList, k.Name)
	}
	return
}

//GetDeployment - get the deployment details for the given app.deployment
func (s *ApplicationDeployments) GetDeployment(appName, deploymentName string) (deployment Deployment) {
	for _, d := range s.Applications[appName].Deployments {

		if d.Name == deploymentName {
			deployment = d
		}
	}
	return
}
