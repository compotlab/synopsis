var app = angular.module("PackageApp", ["ngResource"]);

app.factory("Packages", function ($resource) {
    return $resource("/package/all");
});

app.controller("AllPackages", function ($scope, Packages) {
    Packages.get(function (data) {
        $scope.packages = data.packages;
    });
});

app.filter("Find", function() {
    return function(input, str) {
        var tmp = {};
        angular.forEach(input, function(val, key) {
            if (key.indexOf(str) !== -1) {
                tmp[key] = val;
            }
        });
        return tmp;
    }
});
