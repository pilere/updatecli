package pipeline

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/updatecli/updatecli/pkg/core/result"
)

// RunSources iterates on every source definition to retrieve every information.
func (p *Pipeline) RunSources() error {

	logrus.Infof("\n\n%s\n", strings.ToTitle("Sources"))
	logrus.Infof("%s\n", strings.Repeat("=", len("Source")+1))

	sortedSourcesKeys, err := SortedSourcesKeys(&p.Sources)
	if err != nil {
		logrus.Errorf("%s %v\n", result.FAILURE, err)
		return err
	}

	for _, id := range sortedSourcesKeys {
		err = p.Config.Update(p)
		if err != nil {
			return err
		}

		source := p.Sources[id]
		source.Config = p.Config.Sources[id]

		rpt := p.Report.Sources[id]

		logrus.Infof("\n%s\n", id)
		logrus.Infof("%s\n", strings.Repeat("-", len(id)))

		err = source.Run()
		rpt.Result = source.Result

		if len(source.Changelog) > 0 {
			logrus.Infof("\n\n%s:\n", strings.ToTitle("Changelog"))
			logrus.Infof("%s\n", strings.Repeat("-", len("Changelog")+1))
			logrus.Infof("%s\n", source.Changelog)
		}

		if err != nil {
			logrus.Errorf("%s %v\n", source.Result, err)
		}

		if strings.Compare(source.Result, result.ATTENTION) == 0 {
			logrus.Infof("\n%s empty source returned", source.Result)
		}

		p.Sources[id] = source
		p.Report.Sources[id] = rpt

	}

	return err
}
