/* jslint node: true */
"use strict";

function NoteCtrl($scope, $stateParams, note) {
  var vm = this;

  vm.note = note;

  $scope.$watch('vm.note.content', function(current, original) {
    vm.outputText = marked(current);
  });
}

NoteCtrl.resolve = {
  note: function($stateParams, ApiService, $rootScope) {
    let uuid = $stateParams.uuid;

    return ApiService.getNote(uuid).then(function(result) {
      let note = result;
      note.edit = true;
      $rootScope.$broadcast('currentNote', note);

      return note;
    });
  }
};

angular.module("app").controller("NoteCtrl", NoteCtrl);
