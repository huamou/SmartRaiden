package restful2

import (
	"github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/restful2/v1"
)

func init() {

}

func Start(RaidenApi *smartraiden.RaidenApi, config *params.Config) {
	v1.RaidenApi = RaidenApi
	v1.Config = config
	v1.Start()
}
