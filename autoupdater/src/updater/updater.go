package updater

type Updater interface {
	Update(service string)
}

func NewSwarm() Updater{
	return newSwarmUpdater()
}