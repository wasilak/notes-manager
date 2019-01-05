/* jslint node: true */
"use strict";

angular.module("app").component("list", 
  {
    controller: function ListCtrl($rootScope, ApiService, $scope) {
      var vm = this;

      vm.list = {
        success: true
      };

      vm.sort = "updated:desc";

      vm.tags = [];

      $rootScope.$on('currentNote', function(event, note) {
        vm.note = note;
      });

      vm.updateList = function() {
        ApiService.getList(vm.listFilter, vm.sort, vm.tags).then(function(result) {
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

      vm.loadItems = function(query) {
        return ApiService.getTags(query);
      };

      $scope.$watch('$ctrl.tags', function(current, original) {
        vm.updateList();
      });

      vm.clearSearch();

      vm.updateList();
    },
    templateUrl: "/static/app/views/list.html"
});
