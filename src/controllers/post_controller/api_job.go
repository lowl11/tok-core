package post_controller

import (
	"github.com/lowl11/lazylog/layers"
	"tok-core/src/definition"
)

/*
	FillUnauthorizedJob обертка для _fillUnauthorizedFeed
*/
func (controller *Controller) FillUnauthorizedJob() error {
	if err := controller._fillUnauthorizedFeed(); err != nil {
		return err
	}

	definition.Logger.Info("Fill unauthorized job done!", layers.Job)

	return nil
}

/*
	FillExploreJob обертка для _fillExploreFeed
*/
func (controller *Controller) FillExploreJob() error {
	if err := controller._fillExploreFeed(); err != nil {
		return err
	}

	definition.Logger.Info("Fill explore job done!", layers.Job)

	return nil
}
