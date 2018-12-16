/* jslint node: true */
"use strict";

angular.module("app").component("note", 
  {
    bindings: {
      note: '<'
    },
    controller: function ($scope, $rootScope, $stateParams, ApiService, $state) {
      var vm = this;

      vm.saveNote = function() {
        ApiService.saveNote(vm.note).then(function(result) {
          // some kind of message, i.e. growl
          vm.note = result;
          $rootScope.notifications.push('Note saved');
          $state.go('list_note', {uuid: vm.note.response.id});
        });
      };

      vm.deleteNote = function() {
        var confirmed = confirm("Are you sure?");

        if (confirmed) {
          ApiService.deleteNote(vm.note.response.id).then(function(result) {
              $rootScope.notifications.push('Note deleted');
              $state.go('list', {}, {reload: true});
          });
        }
      };

      $scope.$watch('$ctrl.note.response.content', function(current, original) {
        vm.errorMessage = false;
        try {
          vm.outputText = marked(current);
        } catch (err) {
          vm.errorMessage = err.message;
        }
      });
    },
    templateUrl: "/static/app/views/note.html"
  }
);
