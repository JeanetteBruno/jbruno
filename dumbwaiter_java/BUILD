# note: in order to use this BUILD file, dumbwaiter_java needs to be copied
# to ~/workspace/source/sandbox/users/jbruno/testing/java/dumbwaiter_java
java_library(
  name='controller',
  dependencies = [
    '3rdparty/jvm/junit:junit',
  ],
  sources=['src/main/java/**/*.java']
  )

junit_tests(
  name='controller_tests',
  dependencies = [
    'sandbox/users/jbruno/testing/java/dumbwaiter_java:controller',
  ],
  sources=['src/test/java/**/*.java']
)
