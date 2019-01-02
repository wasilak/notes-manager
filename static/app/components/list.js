/* jslint node: true */
"use strict";

angular.module("app").component("list", 
  {
    controller: function ListCtrl($rootScope, ApiService) {
      var vm = this;

      vm.list = {
        success: true
      };

      vm.sort = "updated:desc";

      $rootScope.$on('currentNote', function(event, note) {
        vm.note = note;
      });

      vm.updateList = function() {
        ApiService.getList(vm.listFilter, vm.sort).then(function(result) {
          vm.list = result;
        });
      };

      vm.setSort = function() {
        vm.updateList();
      };

      vm.clearSearch = function() {
        vm.listFilter = "";
        vm.search();
      };

      vm.search = function() {
        if (vm.listFilter == "") {
          vm.sort = "updated:desc";
        } else {
          vm.sort = "";
        }
        vm.updateList();
      };

      vm.clearSearch();

      vm.updateList();
    },
    templateUrl: "/static/app/views/list.html"
});
