/* jslint node: true */
"use strict";

var app = angular.module("app", ['ui.router', 'ngSanitize', 'growlNotifications'])

.config([function() {
  marked.setOptions({
      renderer: new marked.Renderer(),
      gfm: true,
      tables: true,
      breaks: true,
      pedantic: false,
      sanitize: false, // if false -> allow plain old HTML ;)
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
  $rootScope.user = false;
})
;
