/* jslint node: true */
"use strict";

angular.module("app").component("menu", 
  {
    bindings: {
      // note: '<'
    },
    controller: function($rootScope, ApiService, $state, $scope) {
      var vm = this;

      vm.note = null;

      $rootScope.$on('currentNote', function(event, note) {
        vm.note = note;
      });

      vm.saveNote = function() {
        ApiService.saveNote(vm.note).then(function(result) {
          // some kind of message, i.e. growl
          vm.note = result;
          $rootScope.notifications.push('Note saved');
          $state.go('list_note', {uuid: vm.note.response.id});
        });
      };

      vm.cancel = function(uuid) {
        vm.note = null;

        if (uuid) {
          $state.go('list_note', {uuid: uuid});
        } else {
          $state.go('list', {}, {reload: true});
        }
      };

      vm.createNote = function() {
        ApiService.createNote(vm.note).then(function(result) {
          vm.note = result;
          $rootScope.notifications.push('Note created');
          $state.go('list_note', {uuid: vm.note.response.id}, {reload: true});
        });
      };

      vm.deleteNote = function() {
        var confirmed = confirm("Are you sure?");

        if (confirmed) {
          ApiService.deleteNote(vm.note.response.id).then(function(result) {
              vm.note = null;
              $rootScope.notifications.push('Note deleted');
              $state.go('list', {}, {reload: true});
          });
        }
      };

    },
    templateUrl: "/static/app/views/menu.html"
  }
);