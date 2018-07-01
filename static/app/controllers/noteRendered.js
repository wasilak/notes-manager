/* jslint node: true */
"use strict";

function NoteRenderedCtrl($rootScope, $scope, $stateParams, note) {
  var vm = this;

  vm.uuid = $stateParams.uuid;

  vm.inputText = '';
  vm.note = note;
  $rootScope.$broadcast('currentNote', note);

  vm.outputText = marked(vm.note.content);
}

NoteRenderedCtrl.resolve = {
  note: function($stateParams, ApiService, $rootScope) {
    var uuid = $stateParams.uuid;

    return ApiService.getNote(uuid).then(function(result) {
      $rootScope.$broadcast('currentNote', result);
      return result;
    });
  }
};

angular.module("app").controller("NoteRenderedCtrl", NoteRenderedCtrl);
