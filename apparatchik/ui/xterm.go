package ui

import (
	"fmt"

	"github.com/netice9/apparatchik/apparatchik/core"
	"gitlab.netice9.com/dragan/go-reactor"
	brc "gitlab.netice9.com/dragan/go-reactor/core"
)

type XTerm struct {
	ctx  reactor.ScreenContext
	goal *core.Goal
}

func XTermFactory(ctx reactor.ScreenContext) reactor.Screen {

	appName := ctx.Params["application"]
	goalName := ctx.Params["goal"]

	app, err := core.ApparatchikInstance.GetApplicationByName(appName)
	if err != nil {
		reactor.DefaultNotFoundScreenFactory(ctx)
	}

	goal, found := app.Goals[goalName]
	if !found {
		reactor.DefaultNotFoundScreenFactory(ctx)
	}

	return &XTerm{
		ctx:  ctx,
		goal: goal,
	}
}

var xtermView = brc.MustParseDisplayModel(`
<bs.Panel header="Exec Terminal Session">
	<div id="container" data-api-path="/test" htmlID="terminal-container"></div>
</bs.Panel>
`)

func (x *XTerm) render() {
	view := xtermView.DeepCopy()

	path := fmt.Sprintf("/api/v1.0/applications/%s/goals/%s/exec", x.goal.ApplicationName, x.goal.Name)

	view.SetElementAttribute("container", "data-api-path", path)

	x.ctx.UpdateScreen(&brc.DisplayUpdate{
		Model: WithNavigation(view, [][]string{
			{"Applications", "#/"},
			{x.goal.ApplicationName, fmt.Sprintf("#/apps/%s", x.goal.ApplicationName)},
			{x.goal.Name, fmt.Sprintf("#/apps/%s/%s", x.goal.ApplicationName, x.goal.Name)},
			{"XTerm", fmt.Sprintf("#/apps/%s/%s/xterm", x.goal.ApplicationName, x.goal.Name)},
		}),
	})

	x.ctx.UpdateScreen(&brc.DisplayUpdate{
		Eval: `
		startTerminal()
		`,
	})
}

func (x *XTerm) OnUserEvent(evt *brc.UserEvent) {

}
func (x *XTerm) Mount() {
	x.render()
}

func (x *XTerm) Unmount() {
}
