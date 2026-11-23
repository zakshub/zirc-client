package desktop

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/zakshub/zirc-client/internal/irc"
)

type workspace struct {
	window     fyne.Window
	client     *irc.Client
	config     irc.Config
	connected  bool
	active     string
	channels   []string
	transcript map[string][]string
	members    map[string][]string
	topics     map[string]string

	navigation *widget.List
	memberList *widget.List
	log        *widget.Label
	topic      *widget.Label
	composer   *widget.Entry
	status     *widget.Label
	navItems   []string
	memberRows []string
	disconnectBtn *widget.Button
}

func New(w fyne.Window) *workspace {
	d := &workspace{
		window:     w,
		active:     "Status",
		channels:   []string{"Status"},
		transcript: map[string][]string{"Status": {"ZIRC desktop workspace ready."}},
		members:    make(map[string][]string),
		topics:     make(map[string]string),
	}
	d.build()
	return d
}

func (d *workspace) Show() {
	d.window.Show()
}

func (d *workspace) build() {
	d.window.SetMainMenu(d.menu())

	title := widget.NewLabelWithStyle("ZIRC", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	title.TextStyle = fyne.TextStyle{Bold: true}
	subtitle := widget.NewLabel("Classic IRC workspace")
	subtitle.Importance = widget.LowImportance
	connect := widget.NewButton("Connect", d.showConnect)
	d.disconnectBtn = widget.NewButton("Disconnect", d.disconnect)
	d.disconnectBtn.Disable()
	d.status = widget.NewLabel("Disconnected")
	d.status.Importance = widget.LowImportance

	header := container.NewBorder(nil, nil, container.NewHBox(title, subtitle), container.NewHBox(d.status, connect, d.disconnectBtn), widget.NewSeparator())

	d.navItems = append([]string(nil), d.channels...)
	d.navigation = widget.NewList(
		func() int { return len(d.navItems) },
		func() fyne.CanvasObject { return widget.NewLabel("Channel") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(d.navItems[i]) },
	)
	d.navigation.OnSelected = func(i widget.ListItemID) {
		if i >= 0 && i < len(d.navItems) {
			d.activate(d.navItems[i])
		}
	}
	sideTitle := widget.NewLabelWithStyle("NETWORK / CHANNELS", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	sidebar := container.NewBorder(container.NewVBox(sideTitle, widget.NewSeparator()), nil, nil, nil, d.navigation)

	d.topic = widget.NewLabel("Status window")
	d.topic.Wrapping = fyne.TextWrapWord
	d.log = widget.NewLabel("")
	d.log.Wrapping = fyne.TextWrapWord
	d.log.TextStyle = fyne.TextStyle{Monospace: true}
	logScroll := container.NewVScroll(d.log)

	d.composer = widget.NewEntry()
	d.composer.SetPlaceHolder("Type a message or IRC command, then press Enter")
	d.composer.OnSubmitted = func(value string) {
		d.submit(value)
	}
	send := widget.NewButton("Send", func() {
		d.submit(d.composer.Text)
	})
	input := container.NewBorder(nil, nil, nil, send, d.composer)
	chat := container.NewBorder(container.NewVBox(d.topic, widget.NewSeparator()), input, nil, nil, logScroll)

	d.memberRows = []string{}
	d.memberList = widget.NewList(
		func() int { return len(d.memberRows) },
		func() fyne.CanvasObject { return widget.NewLabel("nickname") },
		func(i widget.ListItemID, o fyne.CanvasObject) { o.(*widget.Label).SetText(d.memberRows[i]) },
	)
	memberTitle := widget.NewLabelWithStyle("NICKNAMES", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	members := container.NewBorder(container.NewVBox(memberTitle, widget.NewSeparator()), nil, nil, nil, d.memberList)

	right := container.NewHSplit(chat, members)
	right.SetOffset(0.80)
	main := container.NewHSplit(sidebar, right)
	main.SetOffset(0.21)

	d.window.SetContent(container.NewBorder(header, nil, nil, nil, main))
	d.refresh()
}

func (d *workspace) menu() *fyne.MainMenu {
	connect := fyne.NewMenuItem("Connect...", d.showConnect)
	disconnect := fyne.NewMenuItem("Disconnect", d.disconnect)
	quit := fyne.NewMenuItem("Quit", func() { d.disconnect(); d.window.Close() })
	file := fyne.NewMenu("File", connect, disconnect, quit)
	join := fyne.NewMenuItem("Join channel...", d.showJoin)
	leave := fyne.NewMenuItem("Leave active channel", func() { d.submit("/part") })
	whois := fyne.NewMenuItem("Whois active nickname", func() {
		if d.active != "Status" { d.submit("/whois " + d.active) }
	})
	channel := fyne.NewMenu("IRC", join, leave, whois)
	about := fyne.NewMenu("Help", fyne.NewMenuItem("About ZIRC", func() {
		dialog.ShowInformation("About ZIRC", "ZIRC is a programmable IRC client in development. The desktop workspace is being built on the shared Go IRC engine.", d.window)
	}))
	return fyne.NewMainMenu(file, channel, about)
}

func (d *workspace) showConnect() {
	server := widget.NewEntry()
	server.SetText("irc.libera.chat:6697")
	nick := widget.NewEntry()
	nick.SetText("zirc-user")
	username := widget.NewEntry()
	username.SetText("zirc-user")
	realname := widget.NewEntry()
	realname.SetText("ZIRC User")
	channels := widget.NewEntry()
	channels.SetPlaceHolder("#channel, #another")
	tlsEnabled := widget.NewCheck("Use TLS with certificate verification", nil)
	tlsEnabled.SetChecked(true)
	items := []*widget.FormItem{
		widget.NewFormItem("Server host:port", server),
		widget.NewFormItem("Nickname", nick),
		widget.NewFormItem("Username", username),
		widget.NewFormItem("Real name", realname),
		widget.NewFormItem("Auto join", channels),
		widget.NewFormItem("Security", tlsEnabled),
	}
	dialog.ShowForm("Connect to IRC network", "Connect", "Cancel", items, func(ok bool) {
		if !ok { return }
		host := strings.TrimSpace(server.Text)
		n := strings.TrimSpace(nick.Text)
		if host == "" || n == "" {
			dialog.ShowError(fmt.Errorf("server and nickname are required"), d.window)
			return
		}
		user := strings.TrimSpace(username.Text)
		if user == "" { user = n }
		real := strings.TrimSpace(realname.Text)
		if real == "" { real = "ZIRC User" }
		var join []string
		for _, ch := range strings.Split(channels.Text, ",") {
			ch = strings.TrimSpace(ch)
			if ch != "" { join = append(join, ch) }
		}
		d.connect(irc.Config{Address: host, Nick: n, Username: user, RealName: real, TLS: tlsEnabled.Checked, Channels: join})
	}, d.window)
}

func (d *workspace) showJoin() {
	channel := widget.NewEntry()
	channel.SetPlaceHolder("#channel")
	dialog.ShowForm("Join channel", "Join", "Cancel", []*widget.FormItem{widget.NewFormItem("Channel", channel)}, func(ok bool) {
		if ok && strings.TrimSpace(channel.Text) != "" { d.submit("/join " + strings.TrimSpace(channel.Text)) }
	}, d.window)
}

func (d *workspace) connect(cfg irc.Config) {
	if d.connected {
		dialog.ShowInformation("Already connected", "Disconnect the current network before connecting another one. Multiple simultaneous networks are part of a later milestone.", d.window)
		return
	}
	d.status.SetText("Connecting to " + cfg.Address + "...")
	go func() {
		client := irc.NewClient(cfg)
		if err := client.Connect(); err != nil {
			fyne.Do(func() {
				d.status.SetText("Connection failed")
				dialog.ShowError(err, d.window)
			})
			return
		}
		fyne.Do(func() {
			d.client = client
			d.config = cfg
			d.connected = true
			d.disconnectBtn.Enable()
			d.status.SetText("Connected to " + cfg.Address)
			for _, ch := range cfg.Channels { d.addChannel(ch) }
			d.append("Status", fmt.Sprintf("* Connected to %s as %s", cfg.Address, cfg.Nick))
			d.refresh()
		})
		go d.readEvents(client)
	}()
}

func (d *workspace) readEvents(client *irc.Client) {
	events, done := client.Events(), client.Done()
	for {
		select {
		case event := <-events:
			fyne.Do(func() { d.handle(event) })
		case <-done:
			fyne.Do(func() {
				if d.client == client {
					d.connected = false
					d.disconnectBtn.Disable()
					d.status.SetText("Disconnected")
					d.append("Status", "* Disconnected from " + d.config.Address)
					d.refresh()
				}
			})
			return
		}
	}
}

func (d *workspace) disconnect() {
	if d.client != nil {
		client := d.client
		d.client = nil
		d.connected = false
		if d.disconnectBtn != nil { d.disconnectBtn.Disable() }
		client.Close()
		d.status.SetText("Disconnected")
		d.append("Status", "* Disconnected")
		d.refresh()
	}
}

func (d *workspace) handle(e irc.Event) {
	switch e.Kind {
	case "status":
		d.append("Status", "* "+e.Text)
		d.status.SetText(e.Text)
	case "error":
		d.append(d.active, "! "+e.Text)
		d.status.SetText("Network error")
	case "message":
		target := e.Target
		if !strings.HasPrefix(target, "#") && !strings.HasPrefix(target, "&") {
			target = e.Nick
		}
		d.addChannel(target)
		d.append(target, "<"+e.Nick+"> "+e.Text)
	case "join":
		d.addChannel(e.Target)
		d.members[e.Target] = appendUnique(d.members[e.Target], e.Nick)
		d.append(e.Target, "* "+e.Nick+" joined")
	case "part":
		d.members[e.Target] = remove(d.members[e.Target], e.Nick)
		d.append(e.Target, "* "+e.Nick+" left")
	case "quit":
		for ch := range d.members { d.members[ch] = remove(d.members[ch], e.Nick) }
		d.append(d.active, "* "+e.Nick+" quit")
	case "names":
		d.addChannel(e.Target)
		for _, name := range strings.Fields(e.Text) {
			name = strings.TrimLeft(name, "~&@%+")
			if name != "" { d.members[e.Target] = appendUnique(d.members[e.Target], name) }
		}
	case "topic":
		d.addChannel(e.Target)
		d.topics[e.Target] = e.Text
		d.append(e.Target, "* Topic: "+e.Text)
	case "notice":
		d.append("Status", e.Text)
	}
	d.refresh()
}

func (d *workspace) submit(value string) {
	value = strings.TrimSpace(value)
	if value == "" { return }
	d.composer.SetText("")
	if strings.HasPrefix(value, "/") {
		fields := strings.Fields(value)
		switch strings.ToLower(fields[0]) {
		case "/join":
			if len(fields) > 1 && d.client != nil {
				d.client.Send("JOIN", fields[1])
				d.addChannel(fields[1])
				d.activate(fields[1])
			}
		case "/part":
			if d.client != nil && d.active != "Status" {
				d.client.Send("PART", d.active)
				d.append(d.active, "* Leaving "+d.active)
				d.channels = remove(d.channels, d.active)
				d.activate("Status")
			}
		case "/nick":
			if len(fields) > 1 && d.client != nil { d.client.Send("NICK", fields[1]); d.config.Nick = fields[1] }
		case "/msg", "/query":
			if len(fields) > 1 {
				target := fields[1]
				d.addChannel(target)
				d.activate(target)
				if len(fields) > 2 && d.client != nil {
					body := strings.TrimSpace(strings.TrimPrefix(value, fields[0]+" "+target))
					d.client.Send("PRIVMSG", target, body)
					d.append(target, "<you> "+body)
				}
			}
		case "/me":
			if d.client != nil && d.active != "Status" && len(fields) > 1 {
				body := strings.TrimSpace(strings.TrimPrefix(value, fields[0]))
				d.client.Send("PRIVMSG", d.active, "\x01ACTION "+body+"\x01")
				d.append(d.active, "* "+d.config.Nick+" "+body)
			}
		case "/whois":
			if len(fields) > 1 && d.client != nil { d.client.Send("WHOIS", fields[1]) }
		case "/who":
			target := d.active
			if len(fields) > 1 { target = fields[1] }
			if d.client != nil { d.client.Send("WHO", target) }
		case "/list":
			if d.client != nil { d.client.Send("LIST") }
		case "/away":
			if d.client != nil { d.client.Send("AWAY", strings.TrimSpace(strings.TrimPrefix(value, fields[0]))) }
		case "/topic":
			if d.client != nil && d.active != "Status" {
				if len(fields) > 1 { d.client.Send("TOPIC", d.active, strings.TrimSpace(strings.TrimPrefix(value, fields[0]))) } else { d.client.Send("TOPIC", d.active) }
			}
		case "/invite":
			if d.client != nil && d.active != "Status" && len(fields) > 1 { d.client.Send("INVITE", fields[1], d.active) }
		case "/kick", "/ban", "/unban", "/op", "/deop", "/voice", "/devoice", "/mode":
			d.moderate(fields)
		case "/notice":
			if d.client != nil && len(fields) > 2 {
				body := strings.TrimSpace(strings.TrimPrefix(value, fields[0]+" "+fields[1]))
				d.client.Send("NOTICE", fields[1], body)
			}
		case "/clear":
			d.transcript[d.active] = nil
			d.refresh()
		case "/quit":
			d.disconnect()
		case "/help":
			d.append("Status", "Commands: /join /part /nick /msg /query /me /whois /who /list /away /topic /invite /kick /ban /unban /op /deop /voice /devoice /mode /notice /clear /quit")
			d.activate("Status")
		default:
			d.append(d.active, "! Unknown desktop command: "+fields[0])
		}
		return
	}
	if d.client == nil || !d.connected {
		d.append("Status", "! Connect to an IRC network before sending messages")
		d.activate("Status")
		return
	}
	if d.active == "Status" {
		d.append("Status", "! Select a channel or open a private conversation first")
		return
	}
	d.client.Send("PRIVMSG", d.active, value)
	d.append(d.active, "<"+d.config.Nick+"> "+value)
	d.refresh()
}

func (d *workspace) moderate(fields []string) {
	if d.client == nil || d.active == "Status" { return }
	if len(fields) < 2 { d.append(d.active, "! Missing command argument"); return }
	switch strings.ToLower(fields[0]) {
	case "/kick":
		reason := "Removed by channel operator"
		if len(fields) > 2 { reason = strings.TrimSpace(strings.TrimPrefix(strings.Join(fields, " "), fields[0]+" "+fields[1])) }
		d.client.Send("KICK", d.active, fields[1], reason)
	case "/ban": d.client.Send("MODE", d.active, "+b", fields[1])
	case "/unban": d.client.Send("MODE", d.active, "-b", fields[1])
	case "/op": d.client.Send("MODE", d.active, "+o", fields[1])
	case "/deop": d.client.Send("MODE", d.active, "-o", fields[1])
	case "/voice": d.client.Send("MODE", d.active, "+v", fields[1])
	case "/devoice": d.client.Send("MODE", d.active, "-v", fields[1])
	case "/mode": d.client.Send("MODE", append([]string{d.active}, fields[1:]...)...)
	}
}

func (d *workspace) activate(target string) {
	d.active = target
	d.refresh()
}

func (d *workspace) addChannel(ch string) {
	if ch == "" { return }
	for _, existing := range d.channels { if existing == ch { return } }
	d.channels = append(d.channels, ch)
	d.transcript[ch] = []string{}
	d.members[ch] = []string{}
}

func (d *workspace) append(target, line string) {
	if target == "" { target = "Status" }
	d.transcript[target] = append(d.transcript[target], time.Now().Format("15:04:05")+" "+line)
	if len(d.transcript[target]) > 3000 { d.transcript[target] = d.transcript[target][len(d.transcript[target])-3000:] }
}

func (d *workspace) refresh() {
	d.navItems = append([]string(nil), d.channels...)
	d.navigation.Refresh()
	d.log.SetText(strings.Join(d.transcript[d.active], "\n"))
	d.topic.SetText(d.topics[d.active])
	if d.active == "Status" { d.topic.SetText("Server status and network events") }
	d.memberRows = append([]string(nil), d.members[d.active]...)
	d.memberList.Refresh()
}

func appendUnique(items []string, value string) []string {
	for _, item := range items { if item == value { return items } }
	return append(items, value)
}

func remove(items []string, value string) []string {
	out := items[:0]
	for _, item := range items { if item != value { out = append(out, item) } }
	return out
}
