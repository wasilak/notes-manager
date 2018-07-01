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

  vm.cancel = function(uuid) {
    vm.note = null;
    $rootScope.$broadcast('currentNote', vm.note);

    if (uuid) {
      $state.go('parent.list.note', {uuid: uuid});
    } else {
      $state.go('parent.list', {}, {reload: true});
    }
  };

  vm.createNote = function() {
    ApiService.createNote(vm.note).then(function(result) {
      vm.note = result;
      $rootScope.notifications.push('Note created');
      $rootScope.$broadcast('currentNote', vm.note);
      $state.go('parent.list.note', {uuid: vm.note.id}, {reload: true});
    });
  };

  vm.deleteNote = function() {
    var confirmed = confirm("Are you sure?");

    if (confirmed) {
      ApiService.deleteNote(vm.note.id).then(function(result) {
          vm.note = null;
          $rootScope.$broadcast('currentNote', vm.note);
          $rootScope.notifications.push('Note deleted');
          $state.go('parent.list', {}, {reload: true});
      });
    }
  };
}

MenuCtrl.resolve = {
};

angular.module("app").controller("MenuCtrl", MenuCtrl);
