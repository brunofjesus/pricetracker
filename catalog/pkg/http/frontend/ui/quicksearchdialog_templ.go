// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.513
package ui

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func QuickSearchDialogComponent() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"modal modal-lg\" id=\"quickSearchDialogModal\"><div class=\"modal-dialog modal-dialog-scrollable\" role=\"document\"><div class=\"modal-content\"><div class=\"modal-body px-3 py-3\"><div class=\"d-flex flex-row align-items-center gap-1\" hx-ext=\"client-side-templates\"><input id=\"quick_search_input\" type=\"search\" class=\"form-control\" placeholder=\"Search...\" aria-label=\"Search\" name=\"name\" hx-get=\"/api/v1/product/search\" hx-trigger=\"keyup changed delay:500ms\" hx-target=\"#quick_search_results\" handlebars-template=\"quick_search_template\"> <button type=\"button\" class=\"btn btn-outline-dark\" data-bs-dismiss=\"modal\" aria-label=\"Close\"><i class=\"bi bi-x-lg\"></i></button></div><div class=\"row mt-3\"><div class=\"list-group\" id=\"quick_search_results\"></div></div></div></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.Raw(`
    <template id="quick_search_template">
    {{#items}}
       <a href="/{{product_id}}" class="list-group-item list-group-item-action">
           <div class="row">
               <div class="col-md-2">
                 <img class="d-block" src="{{image_url}}" width="80px"/>
               </div>
               <div class="col-md-8">
                 <strong class="list-group-item-heading">{{name}}</strong>
                 <p class="list-group-item-text">{{store_name}}</p>
               </div>
               <div class="col-md-2">
                  {{#when discount_percent 'gt' 0}}
                      <span class="badge alert-success">{{math price '/' 100}} €</span>
                  {{else}}
                      <span class="badge alert-danger">{{math price '/' 100}} €</span>
                  {{/when}}
               </div>
           </div>
       </a>
    {{/items}}
    </template>
    `).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
