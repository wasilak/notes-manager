var codeMirror = function($timeout){
  return {
    restrict: "E",
    replace: true,
    require: "?ngModel",
    transclude: true,
    scope: {
      syntax: "@",
      theme: "@",
      ngModel: "="
    },
    template: '<div class="code-editor"></div>',
    link: function(scope, element, attrs, ngModelCtrl, transclude){
      var editor = CodeMirror(element[0], {
        mode: scope.syntax || "javascript",
        theme: scope.theme || "default",
        autoCloseBrackets: scope.autoCloseBrackets || true,
        matchBrackets: scope.matchBrackets || true,
        lineNumbers: true
      });

      if(ngModelCtrl) {
        $timeout(function(){
          ngModelCtrl.$render = function() {
            editor.setValue(ngModelCtrl.$viewValue);
          };
        });
      }

      transclude(function(clonedEl){
//            var initialText = clonedEl.text();
        var initialText = scope.ngModel;
        editor.setValue(initialText);

        if(ngModelCtrl){
          $timeout(function(){
            if(initialText && !ngModelCtrl.$viewValue){
              ngModelCtrl.$setViewValue(initialText);
            }

            editor.on('change', function(){
              ngModelCtrl.$setViewValue(editor.getValue());
            });
          });
        }
      });

      scope.$on('$destroy', function(){
        editor.off('change');
      });
    }
  };
};

angular.module("app").directive("codeEditor", codeMirror);
