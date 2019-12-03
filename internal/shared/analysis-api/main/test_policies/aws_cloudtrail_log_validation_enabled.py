# Tags: ['CIS', 'AWS Managed Rules - Management and Governance']
def policy(resource):
    # Explicit check for True as the value may be None, and we want to return a bool not a NoneType
    return resource['Info']['LogFileValidationEnabled'] is True
