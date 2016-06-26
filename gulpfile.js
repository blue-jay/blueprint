var gulp = require('gulp');

// Style Task
gulp.task('style', function() {
    var sass = require('gulp-sass');
    gulp.src('asset/dynamic/sass/**/*.scss')
        // Available for outputStyle: expanded, nested, compact, compressed
        .pipe(sass({outputStyle: 'expanded'}).on('error', sass.logError))
        .pipe(gulp.dest('asset/static/css/'));
});

// Watch
gulp.task('watch', function() {
    gulp.watch('asset/dynamic/sass/**/*.scss', ['style']);
});

// Default
gulp.task('default', ['style']);