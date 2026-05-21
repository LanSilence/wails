//go:build android

package runtime

// Android uses window.wails.invoke which is set up via addJavascriptInterface in WailsJSBridge
var invoke = `
console.log('[Wails Android Runtime] Injecting runtime');
window._wails.invoke=function(m){
    return window.wails.invoke(typeof m==='string'?m:JSON.stringify(m));
};
// Override fetch to intercept /wails/runtime calls and route them through JNI
(function() {
  var origFetch = window.fetch;
  window.fetch = function(url, options) {
    var urlStr = (typeof url === 'string') ? url : (url ? url.href : '');
    if (urlStr.indexOf('/wails/runtime') !== -1) {
      try {
        var bodyStr = (options && options.body) || '{}';
        var body = (typeof bodyStr === 'string') ? JSON.parse(bodyStr) : bodyStr;
        var response = window._wails.invoke(JSON.stringify(body));
        var result = JSON.parse(response);
        return Promise.resolve({
          ok: true, status: 200,
          headers: { get: function() { return 'application/json'; } },
          json: function() { return Promise.resolve(result); },
          text: function() { return Promise.resolve(response); }
        });
      } catch(e) {
        return Promise.reject(e);
      }
    }
    return origFetch.apply(this, arguments);
  };
})();
console.log('[Wails Android Runtime] Runtime injection complete');
`
var flags = ""
