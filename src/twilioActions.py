import src.databaseAccess as databaseAccess
from twilio.rest import Client
import os


# send the submission to the supervisor and myself
def send_text(id: int, msg: str) -> bool:
    phone = databaseAccess.get_employee_phone(id)
    if not phone:
        return False

    client = Client(
        os.environ['TWILIO_ACCOUNT_SID'],
        os.environ['TWILIO_AUTH_TOKEN']
    )

    twilio = os.environ['TWILIO_PHONE']
    client.messages.create(
        from_=f"+1{twilio}",
        to=f"+1{phone}",
        body=msg
    )

    return True

# send confirmation texts of a submission to myself, the supervisor, and the recipient
def confirm_submission(employee_id: int, msg: str, fro: str):
    # text confirmation to owner
    owner_id = databaseAccess.get_employee_id('jr', 'poulsen')
    if not owner_id:
        return "Error. Owner not found..."
    owner_phone = databaseAccess.get_employee_phone(owner_id)
    if not owner_phone:
        return "Error. Owner phone not found."
    send_text(owner_id, msg)

    # text confirmation to supervisor
    supervisor_id = databaseAccess.get_super_id(employee_id)
    if supervisor_id != owner_id: # don't text the supervisor twice if it is also the owner
        if not supervisor_id:
            return "Error. Supervisor not found."
        if supervisor_id > 0: # -1 means no supervisor, so text no one
            supervisor_phone = databaseAccess.get_employee_phone(supervisor_id)
            if not supervisor_phone:
                return "Error. Supervisor phone not found."
            send_text(supervisor_id, msg)

    # text confirmation to admin
    admin_id = databaseAccess.get_employee_id('admin', 'admin')
    if not admin_id:
        return "Error. Admin not found..."
    admin_phone = databaseAccess.get_employee_phone(admin_id)
    if not admin_phone:
        return "Error. Admin phone not found."
    send_text(admin_id, msg)

    # text confirmation to recipient if they aren't already texted in the response
    phone = databaseAccess.get_employee_phone(employee_id)
    if phone != fro.replace("+1", ""):
        send_text(employee_id, msg)

    return True
