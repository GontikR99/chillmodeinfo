// +build wasm,web

package toast

import (
	"github.com/GontikR99/chillmodeinfo/pkg/document"
	"strconv"
	"syscall/js"
	"time"
)

var nextToastId=0
const toastHolderId="toast-holder"

func Error(subsystem string, err error) {
	Popup("Error in "+subsystem, err.Error(), 8*time.Second)
}

func Popup(title string, messageText string, delay time.Duration) {
	toastId := nextToastId
	nextToastId++
	toastIdStr := "error-toast-"+strconv.Itoa(toastId)

	toastDiv := document.CreateElement("div")
	toastDiv.SetAttribute("class", "toast")
	toastDiv.SetAttribute("role", "alert")
	toastDiv.SetAttribute("id", toastIdStr)
	toastDiv.SetAttribute("data-delay", strconv.Itoa(int(delay/time.Millisecond)))

	headerDiv := document.CreateElement("div")
	headerDiv.SetAttribute("class", "toast-header")
	toastDiv.AppendChild(headerDiv)

	titleElt := document.CreateElement("strong")
	titleElt.SetAttribute( "class", "mr-auto")
	headerDiv.AppendChild(titleElt)

	titleText := document.CreateTextNode(title)
	titleElt.AppendChild(titleText)

	closeButton := document.CreateElement("button")
	closeButton.SetAttribute( "type", "button")
	closeButton.SetAttribute( "class", "ml-2 mb-1 close")
	closeButton.SetAttribute( "data-dismiss", "toast")
	headerDiv.AppendChild(closeButton)

	closeSpan := document.CreateElement("span")
	closeButton.AppendChild(closeSpan)

	closeSpanText := document.CreateTextNode("\u00D7")
	closeSpan.AppendChild(closeSpanText)

	bodyDiv := document.CreateElement("div")
	bodyDiv.SetAttribute( "class", "toast-body")
	toastDiv.AppendChild(bodyDiv)

	bodyText := document.CreateTextNode(messageText)
	bodyDiv.AppendChild(bodyText)

	toastHolder := document.GetElementById(toastHolderId)
	toastHolder.AppendChild(toastDiv)

	js.Global().Call("$", "#"+toastIdStr).Call("toast")
	js.Global().Call("$", "#"+toastIdStr).Call("toast", "show")

	disposeFunction:=new(js.Func)
	*disposeFunction=js.FuncOf(func(_ js.Value, _ []js.Value)interface{} {
		disposeFunction.Release()
		js.Global().Call("$", "#"+toastIdStr).Call("toast", "dispose")

		toastElt:= document.GetElementById(toastIdStr)
		toastElt.Remove()
		return nil
	})
	js.Global().Call("$", "#"+toastIdStr).Call("on", "hidden.bs.toast", disposeFunction)
}