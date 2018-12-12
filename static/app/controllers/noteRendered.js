/* jslint node: true */
"use strict";

function NoteRenderedCtrl($rootScope, $stateParams, note) {
  var vm = this;

  vm.inputText = '';
  vm.note = note;
  $rootScope.$broadcast('currentNote', note);

  vm.outputText = marked(vm.note.response.content);
}

NoteRenderedCtrl.resolve = {
  note: function($stateParams, ApiService, $rootScope) {
    return ApiService.getNote($stateParams.uuid);
  }
};

angular.module("app").controller("NoteRenderedCtrl", NoteRenderedCtrl);
