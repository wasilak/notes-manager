/* jslint node: true */
/* jshint -W117 */
"use strict";

angular.module("app").component("note",
  {
    bindings: {
      note: '<'
    },
    controller: function ($scope, $rootScope, $stateParams, ApiService, $state) {
      var vm = this;

      vm.noteOriginal = false;
      vm.loader = false;

      vm.saveNote = function() {
        vm.loader = true;
        ApiService.saveNote(vm.note).then(function(result) {
          // some kind of message, i.e. growl
          vm.note = result;
          $rootScope.notifications.push('Note saved');
          vm.loader = false;
          $state.go('list_note', {uuid: vm.note.response.id});
        });
      };

      vm.saveButtonDisabled = function() {
        return angular.equals(vm.note, vm.noteOriginal);
      };

      vm.deleteNote = function() {
        var confirmed = confirm("Are you sure?");

        if (confirmed) {
          vm.loader = true;
          ApiService.deleteNote(vm.note.response.id).then(function(result) {
              $rootScope.notifications.push('Note deleted');
              vm.loader = false;
              $state.go('list', {}, {reload: true});
          });
        } else {
          vm.loader = false;
        }
      };

      vm.loadItems = function(query) {
        return ApiService.getTags(query);
      };

      vm.breakpoints = [];

      $scope.$watch('$ctrl.note.response', function(current, original) {
        vm.errorMessage = false;
        $rootScope.$state.current.data.title = current.title + " [[edit]]";
        vm.breakpoints = [];
        try {

          // making a copy of original model in order to detect changes and to be able to enable/disable save button
          if (!vm.noteOriginal) {
            vm.noteOriginal = JSON.parse(JSON.stringify(vm.note));
          }
          // vm.outputText = marked(current.content);
          let markdownlintOptions = {
            "config": {
              "MD041": false,
              "MD013": false,
              "MD034": false,
              // "MD031": false,
              // "MD032": false
            },
            "strings": {
              "content": current.content
            }
          }

          let lintResult = markdownlint.sync(markdownlintOptions)

          vm.lintResult = lintResult.toString().replace(/(?:\r\n|\r|\n)/g, '<br>');
          vm.outputText =  marked.parse(current.content);

          for (let id in lintResult.content) {
            vm.breakpoints.push(lintResult.content[id].lineNumber - 1);
          }

        } catch (err) {
          vm.errorMessage = err.message;
        }
      }, true);
    },
    templateUrl: "/static/app/views/note.html"
  }
);
