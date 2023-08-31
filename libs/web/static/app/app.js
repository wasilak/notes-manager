/* jslint node: true */
"use strict";

var app = angular.module("app", ['ui.router', 'ngSanitize', 'growlNotifications', 'ngTagsInput', 'ngclipboard'])

  .config([function () {
    var renderer = new marked.Renderer();

    // opening links in new tab (default link renderer override)
    renderer.link = function (href, title, text) {
      var link = marked.Renderer.prototype.link.call(this, href, title, text);
      return link.replace("<a", "<a target='_blank' ");
    };

    renderer.listitem = function (text, task, checked) {
      if (task && checked) {
        return '<li class="todo_checkbox"><i class="fa fa-square-check" aria-hidden="true"></i>' + text + '</li>\n';
      } else if (task && !checked) {
        return '<li class="todo_checkbox"><i class="fa fa-square" aria-hidden="true"></i>' + text + '</li>\n';
      } else {
        return '<li>' + text + '</li>\n';
      }
    };

    marked.setOptions({
      renderer: renderer,
      gfm: true,
      tables: true,
      breaks: true,
      pedantic: false,
      smartLists: true,
      mangle: false,
      headerIds: false,
      smartypants: false
    });

    marked.use(markedHighlight.markedHighlight({
      langPrefix: 'language-',
      highlight(code, lang) {
        if (lang) {
          return hljs.highlight(code, { language: lang }).value;
        } else {
          return hljs.highlightAuto(code).value;
        }
      }
    }));
  }])

  .config(["$locationProvider", function ($locationProvider) {
    $locationProvider.html5Mode(true);
  }])

  .run(['$rootScope', function ($rootScope) {
    $rootScope.notifications = [];
  }])
  ;
