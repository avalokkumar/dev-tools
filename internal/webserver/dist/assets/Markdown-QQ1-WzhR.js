import{c as b,r,j as e,f as v}from"./index-qUUxopbW.js";import{K as x,k as f,T as m,E as h,a as j,b as w,L as T,S as c}from"./Input-ClLWRu-h.js";import{T as C,u as g}from"./ToolPage-Dvu9Ksnc.js";import{F as S}from"./FormatterShell-9sDB3XZd.js";import{C as u,a as p}from"./Card-BxtjsK-V.js";import{B as k}from"./Button-05bF5-xG.js";import{T as M}from"./Tabs-DyfBae8V.js";import{P as L}from"./play-sQ-ig8C-.js";/**
 * @license lucide-react v0.453.0 - ISC
 *
 * This source code is licensed under the ISC license.
 * See the LICENSE file in the root directory of this source tree.
 */const N=b("Eye",[["path",{d:"M2.062 12.348a1 1 0 0 1 0-.696 10.75 10.75 0 0 1 19.876 0 1 1 0 0 1 0 .696 10.75 10.75 0 0 1-19.876 0",key:"1nclc0"}],["circle",{cx:"12",cy:"12",r:"3",key:"1v7zrd"}]]),E=v("markdown"),P=`# DevForge

A **local-first** developer toolkit.

- 75 operations
- CLI · Web · MCP

\`\`\`go
fmt.Println("hello")
\`\`\`
`,y=`name,role
Alice,Engineer
Bob,Designer
`;function U(){const[t,o]=r.useState("preview");return e.jsxs(C,{tool:E,children:[e.jsx(M,{tabs:[{id:"preview",label:"Markdown → HTML"},{id:"csv",label:"CSV → Table"}],active:t,onChange:l=>o(l),className:"mb-md"}),t==="preview"?e.jsx(A,{}):e.jsx(F,{})]})}function A(){const[t,o]=r.useState(P),[l,d]=r.useState(!0),[i,a]=r.useState(!1),s=g(n=>x(n.input,n.gfm,n.unsafe));return e.jsxs("div",{className:"grid grid-cols-1 lg:grid-cols-2 gap-md",children:[e.jsxs(u,{children:[e.jsx(p,{title:"Markdown"}),e.jsx(f,{rows:14,value:t,onChange:n=>o(n.target.value)}),e.jsxs("div",{className:"grid grid-cols-2 gap-md mt-md",children:[e.jsx(m,{checked:l,onChange:d,label:"GitHub Flavored"}),e.jsx(m,{checked:i,onChange:a,label:"Allow raw HTML",hint:"Unsafe"})]}),e.jsxs(k,{onClick:()=>s.run({input:t,gfm:l,unsafe:i}),loading:s.loading,fullWidth:!0,className:"mt-md",disabled:!t.trim(),children:[e.jsx(L,{className:"h-4 w-4"})," Render"]})]}),e.jsxs(u,{padded:!1,children:[e.jsx(p,{title:"Preview",icon:e.jsx(N,{className:"h-5 w-5"})}),e.jsx("div",{className:"p-md",children:s.error?e.jsx(h,{error:s.error}):s.data?e.jsxs(e.Fragment,{children:[e.jsx("div",{className:"prose-devforge bg-surface-container-lowest border border-outline/10 rounded p-md max-h-[60vh] overflow-auto thin-scrollbar",dangerouslySetInnerHTML:{__html:s.data.output}}),e.jsxs("details",{className:"mt-md",children:[e.jsx("summary",{className:"font-data-label text-data-label uppercase text-on-surface-variant cursor-pointer",children:"View HTML"}),e.jsx(w,{code:s.data.output,language:"html",download:{filename:"preview.html"}})]})]}):e.jsx(j,{title:"Render to preview"})})]})]})}function F(){const[t,o]=r.useState("none"),[l,d]=r.useState(","),i=g(a=>T(a.input,a.delimiter,a.alignment));return e.jsx(S,{initial:y,buildInput:a=>({input:a,delimiter:l,alignment:t}),op:i,language:"markdown",downloadFilename:"table.md",buttonLabel:"Build Table",options:e.jsxs("div",{className:"grid grid-cols-2 gap-md",children:[e.jsx(c,{label:"Alignment",value:t,onChange:a=>o(a.target.value),options:[{value:"none",label:"Default"},{value:"left",label:"Left"},{value:"center",label:"Center"},{value:"right",label:"Right"}]}),e.jsx(c,{label:"Delimiter",value:l,onChange:a=>d(a.target.value),options:[{value:",",label:"Comma"},{value:";",label:"Semicolon"},{value:"	",label:"Tab"}]})]})})}export{U as default};
