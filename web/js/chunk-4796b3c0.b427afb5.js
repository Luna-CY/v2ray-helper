(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-4796b3c0"],{1888:function(e,t,o){},5326:function(e,t,o){"use strict";o.r(t);var n=o("7a23");Object(n["pushScopeId"])("data-v-1ddbc933");var r={class:"login-box"},c={class:"login-form"};function i(e,t,o,i,a,l){var u=Object(n["resolveComponent"])("el-button"),d=Object(n["resolveComponent"])("el-input"),f=Object(n["resolveComponent"])("el-form-item"),s=Object(n["resolveComponent"])("el-form");return Object(n["openBlock"])(),Object(n["createElementBlock"])("div",r,[Object(n["createElementVNode"])("div",c,[Object(n["createVNode"])(s,{ref:"LoginForm",model:e.form,rules:e.rules,"label-width":"0",onsubmit:"return false",onSubmit:e.login},{default:Object(n["withCtx"])((function(){return[Object(n["createVNode"])(f,{prop:"key"},{default:Object(n["withCtx"])((function(){return[Object(n["createVNode"])(d,{modelValue:e.form.key,"onUpdate:modelValue":t[0]||(t[0]=function(t){return e.form.key=t}),placeholder:"口令",size:"medium",autocomplete:"off","prefix-icon":"el-icon-key"},{append:Object(n["withCtx"])((function(){return[Object(n["createVNode"])(u,{loading:e.logging,type:"primary",onClick:e.login,icon:"el-icon-check"},null,8,["loading","onClick"])]})),_:1},8,["modelValue"])]})),_:1})]})),_:1},8,["model","rules","onSubmit"])])])}Object(n["popScopeId"])();var a=o("bc3a"),l=o.n(a),u=o("262e"),d=o("2caf"),f=o("d4ec"),s=o("89df"),m="/api/auth",b=function e(){Object(f["a"])(this,e),this.key=""},p=(s["a"],o("0613")),g=o("6821"),k=Object(n["defineComponent"])({name:"Login",data:function(){return{form:new b,rules:{key:[{required:!0,message:"请输入口令",trigger:"blur"}]},logging:!1}},methods:{login:function(){var e=this,t=this.$refs["LoginForm"];return t.validate((function(t){if(!t)return!1;var o=new b;o.key=g(e.form.key),e.logging=!0,l.a.post(m,o).then((function(t){if(e.logging=!1,0==t.data.code){var o=new p["a"];o.token=t.data.data.token,o.expired=t.data.data.expired,e.$store.commit("token",o);var n=e.$route.query.redirect;window.location.href=n||"/"}else e.$message.error(t.data.message)}))})),!1}}});o("c7fa"),o("f66b");k.render=i,k.__scopeId="data-v-1ddbc933";t["default"]=k},a427:function(e,t,o){},c7fa:function(e,t,o){"use strict";o("1888")},f66b:function(e,t,o){"use strict";o("a427")}}]);
//# sourceMappingURL=chunk-4796b3c0.b427afb5.js.map