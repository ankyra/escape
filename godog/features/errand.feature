Feature: Listing errands
    As an operator
    In order to safely perform day to day tasks
    I need to be able to see what errands I can 

    Scenario: List errands in deployment
      Given a new Escape plan called "my-release"
        And errand "test-errand" with script "test.sh"
        And I release the application
        And I deploy "_/my-release-v0.0.0"
       When I list the errands in the deployment "_/my-release"
       Then I should see "test-errand" in the output

    Scenario: Release should be read from the deployment...
      Given a new Escape plan called "errand-release"
        And errand "test-errand" with script "test.sh"
        And I release the application
      Given a new Escape plan called "errand-release"
        And I release the application
        And I deploy "_/errand-release-v0.0.0"
       When I list the errands in the deployment "_/errand-release"
       Then I should see "test-errand" in the output

    Scenario: ...unless the --local flag is used
      Given a new Escape plan called "my-release"
        And I release the application
        And I deploy "_/my-release-v0.0.0"
      Given a new Escape plan called "my-release"
        And errand "test-errand" with script "test.sh"
       When I list the local errands
       Then I should see "test-errand" in the output
