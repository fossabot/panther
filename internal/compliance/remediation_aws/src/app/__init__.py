from typing import Any, Dict, Type

from .exceptions import RemediationDoesNotExist, RemediationAlreadyExists
from .remediation_base import RemediationBase
from ..remediations import *  # import all files containing subclasses of RemediationBase


class Remediation:
    """Class to be used as a decorator to register all remediation subclasses"""
    _remediations: Dict[str, Type[RemediationBase]] = {}  # Maps remediationID to subclass

    @staticmethod
    def __new__(cls: Type, remediation: Type[RemediationBase]) -> Any:
        remediation_id = remediation.remediation_id()
        cls._remediations[remediation_id] = remediation
        return remediation

    @classmethod
    def get(cls, remediation_id: str) -> Type[RemediationBase]:
        """Return the proper app integration class for this service

        Args:
            remediation_id: The remediation id

        Returns: Subclass of RemediationBase corresponding to the remediation id

        Raises:
            RemediationDoesntExist: Raised if remediation doesn't exist
            for the specified remediation id
        """
        try:
            return cls._remediations[remediation_id]
        except KeyError:
            raise RemediationDoesNotExist('Remediation with id {} does not exist'.format(remediation_id))

    @classmethod
    def get_all_remediations(cls) -> Dict[str, Dict[str, str]]:
        """
        Returns the list of available remediations, along with the parameters
        these remediations use.

        Returns: A map of remediationId -> list of parameters.
        """
        result = {}
        for remediation in cls._remediations.values():
            result[remediation.remediation_id()] = remediation.parameters()
        return result
