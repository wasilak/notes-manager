/* jslint node: true */
"use strict";

var app = angular.module("app", ['ui.router', 'ngSanitize', 'growlNotifications', 'ngTagsInput', 'ngclipboard'])

.config([function() {
  var renderer = new marked.Renderer();
  
  // opening links in new tab (default link renderer override)
  renderer.link = function(href, title, text) {
    var link = marked.Renderer.prototype.link.call(this, href, title, text);
    return link.replace("<a","<a target='_blank' ");
  };

  marked.setOptions({
      renderer: renderer,
      gfm: true,
      tables: true,
      breaks: true,
      pedantic: false,
      sanitizer: DOMPurify.sanitize,
      smartLists: true,
      smartypants: false,
      highlight: function (code, lang) {
        if (lang) {
          return hljs.highlight(lang, code).value;
        } else {
          return hljs.highlightAuto(code).value;
        }
      }
    });
}])

.config(["$locationProvider", function($locationProvider) {
  $locationProvider.html5Mode(true);
}])

.run(function($rootScope) {
  $rootScope.notifications = [];
})
;
