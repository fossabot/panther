"""Policy engine subprocess."""
import json
import sys
from typing import Any, Dict

from .policy import PolicySet


def analyze(data: Dict[str, Any]) -> Dict[str, Any]:
    """Run the Python analysis"""
    policy_set = PolicySet(data['policies'])
    result = {'resources': [policy_set.analyze(r) for r in data['resources']]}
    return result


def main() -> None:
    """Subprocess entry point."""
    process_input = json.loads(sys.stdin.read())
    result = analyze(process_input)

    # Print the json response, which should be faster than going to disk.
    print('\n' + json.dumps(result, separators=(',', ':')))


if __name__ == '__main__':
    main()  # pragma: no cover
