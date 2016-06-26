var gulp = require('gulp');

// SASS Task
gulp.task('sass', function() {
    var sass = require('gulp-sass');
    gulp.src('asset/dynamic/sass/**/*.scss')
        // Available for outputStyle: expanded, nested, compact, compressed
        .pipe(sass({outputStyle: 'expanded'}).on('error', sass.logError))
        .pipe(gulp.dest('asset/static/css/'));
});

// jQuery Task
gulp.task('jquery', function() {
	gulp.src('bower_components/jquery/dist/jquery.min.*')
		.pipe(gulp.dest('asset/static/js/'));
});

// Bootstrap Task
gulp.task('bootstrap', function() {
	gulp.src('bower_components/bootstrap/dist/css/bootstrap-theme.min.*')
		.pipe(gulp.dest('asset/static/css/'));
	gulp.src('bower_components/bootstrap/dist/css/bootstrap.min.*')
		.pipe(gulp.dest('asset/static/css/'));
	gulp.src('bower_components/bootstrap/dist/fonts/*')
		.pipe(gulp.dest('asset/static/fonts/'));
	gulp.src('bower_components/bootstrap/dist/js/bootstrap.min.js')
		.pipe(gulp.dest('asset/static/js/'));
});

// Underscore Task
gulp.task('underscore', function() {
	gulp.src('bower_components/underscore/underscore-min.*')
		.pipe(gulp.dest('asset/static/js/'));
});

// Watch
gulp.task('watch', function() {
    gulp.watch('asset/dynamic/sass/**/*.scss', ['sass']);
});

// Default
gulp.task('default', ['sass', 'jquery', 'bootstrap', 'underscore']);