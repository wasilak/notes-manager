/* jslint node: true */
"use strict";

function NoteCtrl($rootScope, $scope, $stateParams, ApiService, $state) {
  var vm = this;

  vm.uuid = $stateParams.uuid;

  vm.note = {
    content: ''
  };

  ApiService.getNote(vm.uuid).then(function(result) {
    vm.note = result;
    vm.note.edit = true;
    $rootScope.$broadcast('currentNote', vm.note);
  });

  $scope.$watch('vm.note.content', function(current, original) {
    vm.outputText = marked(current);
  });
}

NoteCtrl.resolve = {
};

angular.module("app").controller("NoteCtrl", NoteCtrl);
