/* jslint node: true */
"use strict";

function NoteCtrl($scope, $rootScope, note) {
  var vm = this;

  vm.note = note;
  note.edit = true;
  $rootScope.$broadcast('currentNote', note);

  $scope.$watch('vm.note.content', function(current, original) {
    vm.outputText = marked(current);
  });
}

NoteCtrl.resolve = {
  note: function($stateParams, ApiService, $rootScope) {
    return ApiService.getNote($stateParams.uuid);
  }
};

angular.module("app").controller("NoteCtrl", NoteCtrl);
