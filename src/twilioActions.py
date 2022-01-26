import os

import src.databaseAccess as databaseAccess

from twilio.rest import Client



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
    supervisor_id = databaseAccess.get_super_id(employee_id)
    if supervisor_id != None:
        supervisor_phone = databaseAccess.get_employee_phone(supervisor_id)
        if not supervisor_phone:
            return "Error. Supervisor phone not found."
        send_text(supervisor_id, msg)

    admin_id = databaseAccess.get_employee_id('admin', 'admin')
    if not admin_id:
        return "Error. Admin not found..."
    admin_phone = databaseAccess.get_employee_phone(admin_id)
    if not admin_phone:
        return "Error. Admin phone not found."
    send_text(admin_id, msg)

    return True
