/* jslint node: true */
/* global angular */
/* global marked */
/* global markedHighlight */
/* global hljs */
"use strict";

// eslint-disable-next-line no-unused-vars
var app = angular.module("app", ['ui.router', 'ngSanitize', 'growlNotifications', 'ngTagsInput', 'ngclipboard'])

  .config([function () {
    var renderer = new marked.Renderer();

    // opening links in new tab (default link renderer override)
    renderer.link = function (href, title, text) {
      var link = marked.Renderer.prototype.link.call(this, href, title, text);
      return link.replace("<a", "<a target='_blank' ");
    };

    renderer.listitem = function (item) {

      // original part of function - start
      let itemBody = '';
      if (item.task) {
        const checkbox = this.checkbox({ checked: !!item.checked });
        if (item.loose) {
          if (item.tokens.length > 0 && item.tokens[0].type === 'paragraph') {
            item.tokens[0].text = checkbox + ' ' + item.tokens[0].text;
            if (item.tokens[0].tokens && item.tokens[0].tokens.length > 0 && item.tokens[0].tokens[0].type === 'text') {
              item.tokens[0].tokens[0].text = checkbox + ' ' + item.tokens[0].tokens[0].text;
            }
          }
          else {
            item.tokens.unshift({
              type: 'text',
              raw: checkbox + ' ',
              text: checkbox + ' '
            });
          }
        }
        else {
          itemBody += checkbox + ' ';
        }
      }
      itemBody += this.parser.parse(item.tokens, !!item.loose);
      // original part of function - end

      // my modification - start
      if (item.task) {
        if (item.checked) {
          return `<li class="todo_checkbox"><i class="fa fa-square-check" aria-hidden="true"></i>${itemBody}</li>\n`;
        } else {
          return `<li class="todo_checkbox"><i class="fa fa-square" aria-hidden="true"></i>${itemBody}</li>\n`;
        }
      }
      // my modification - end

      return `<li>${itemBody}</li>\n`;
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
