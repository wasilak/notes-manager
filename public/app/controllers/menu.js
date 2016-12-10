/* jslint node: true */
"use strict";

function MenuCtrl($rootScope, $scope, ApiService, $state) {
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
      $state.go('parent.list.note', {uuid: vm.note.id});
    });
  };

  vm.createNote = function() {
    ApiService.createNote(vm.note).then(function(result) {
      // some kind of message, i.e. growl
      vm.note = result;
      $rootScope.notifications.push('Note created');
      $state.go('parent.list.note', {uuid: vm.note.id});
    });
  };

  vm.deleteNote = function() {
    ApiService.deleteNote(vm.note.id).then(function(result) {
      // some kind of message, i.e. growl

      var confirmed = confirm("Are you sure?");

      if (confirmed) {
        vm.note = null;
        $rootScope.$broadcast('currentNote', vm.note);
        $rootScope.notifications.push('Note deleted');
        $state.go('parent.list', {}, {reload: true});
      }
    });
  };
}

MenuCtrl.resolve = {
};

angular.module("app").controller("MenuCtrl", MenuCtrl);
