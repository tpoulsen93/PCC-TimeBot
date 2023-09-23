import databaseAccess as databaseAccess
from twilio.rest import Client
from enum import Enum
import os

UserType = Enum("UserType", ["Admin", "Owner", "Supervisor", "Employee"])


def send_text(id: int, msg: str):
    phone = databaseAccess.get_employee_phone(id)
    if not phone:
        return

    client = Client(os.environ["TWILIO_ACCOUNT_SID"], os.environ["TWILIO_AUTH_TOKEN"])

    twilio = os.environ["TWILIO_PHONE"]
    client.messages.create(from_=f"+1{twilio}", to=f"+1{phone}", body=msg)


def send_confirmation(id: int, user_type: UserType, is_admin: bool, msg: str):
    if not id:
        return f"Error. {user_type.name} ID not found"
    if id < 1:  # if id < 1, do nothing
        return False

    phone = databaseAccess.get_employee_phone(id)
    if not phone:
        return f"Error. {user_type.name} phone not found"
    if not is_admin:
        send_text(id, msg)
        return False


# send confirmation texts of a submission to myself, the supervisor, and the recipient
def confirm_submission(employee_id: int, msg: str, fro: str):
    # get relevant id's
    admin_id = databaseAccess.get_employee_id("admin", "admin")
    owner_id = databaseAccess.get_employee_id("jr", "poulsen")
    supervisor_id = databaseAccess.get_super_id(employee_id)

    # don't send confirmations to supervisor or owner if the admin submits hours
    is_admin = employee_id == admin_id

    # text confirmation to admin
    admin_confirmation = send_confirmation(admin_id, UserType.Admin, is_admin, msg)
    if admin_confirmation:
        return admin_confirmation

    # text confirmation to owner
    owner_confirmation = send_confirmation(owner_id, UserType.Owner, is_admin, msg)
    if owner_confirmation:
        return owner_confirmation

    # text confirmation to supervisor
    supervisor_confirmation = send_confirmation(
        supervisor_id, UserType.Supervisor, is_admin, msg
    )
    if supervisor_confirmation:
        return supervisor_confirmation

    # text confirmation to recipient if they aren't already texted in the response
    # this is to handle the case that someone submits hours from a phone number other than their usual number.
    # that way, if someone submits another persons hours when they shouldn't, the affected user gets notified
    phone = databaseAccess.get_employee_phone(employee_id)
    if phone != fro.replace("+1", ""):
        send_text(employee_id, msg)

    return True
