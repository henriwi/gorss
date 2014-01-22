var gorss = angular.module('gorss', ["feedService"]);

var feedService = angular.module("feedService", ["ngResource"]);

feedService.factory('Feed', ['$resource',
  function($resource){
    return $resource('api/feeds', {}, {
      query: {method:'GET', isArray:true}
    });
  }]);

gorss.controller("FeedsCtrl", ["$scope", "Feed", function($scope, Feed) {
	$scope.feeds = Feed.query();

	$scope.update = function() {
		$scope.feeds = Feed.query();
	};

}]);