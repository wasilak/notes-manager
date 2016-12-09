/* jslint node: true */
"use strict";

function NoteCtrl($rootScope, $scope, $stateParams, ApiService) {
  var vm = this;

  vm.uuid = $stateParams.uuid;

  vm.note = {
    content: ''
  };

  ApiService.getNote(vm.uuid).then(function(result) {
    vm.note = result;
  });

  $scope.$watch('vm.note.content', function(current, original) {
    vm.outputText = marked(current);
  });

  vm.saveNote = function() {
    ApiService.saveNote(vm.note).then(function(result) {
       // so,e kind of message, i.e. growl
    });
  };
}

NoteCtrl.resolve = {
};

angular.module("app").controller("NoteCtrl", NoteCtrl);
