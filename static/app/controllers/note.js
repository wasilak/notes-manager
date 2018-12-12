/* jslint node: true */
"use strict";

function NoteCtrl($scope, $rootScope, note) {
  var vm = this;

  vm.note = note;
  note.edit = true;
  $rootScope.$broadcast('currentNote', note);

  $scope.$watch('vm.note.response.content', function(current, original) {
    if (!vm.note.success) {
      return;
    }

    vm.errorMessage = false;
    try {
      vm.outputText = marked(current);
    } catch (err) {
      vm.errorMessage = err.message;
    }
  });
}

NoteCtrl.resolve = {
  note: function($stateParams, ApiService, $rootScope) {
    return ApiService.getNote($stateParams.uuid);
  }
};

angular.module("app").controller("NoteCtrl", NoteCtrl);
