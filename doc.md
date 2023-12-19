What is a db transaction?

    A single unit of work
    Often made up of multiple db operations

Example ( money transfer )
Transfer 10 USD from bank account 1 to bank account 2

    1. Create a transfer record with amount = 10
    2. Create an account entry for account 1 with amount = -10
    3. Create an account entry for account 2 with amount = +10
    4. Subtract 10 from the balance of account 1
    5. Add 10 to the balance of account 2
    
Why do we need db transaction?

    1. To provide a reliable and consistent unit of work, even in case of system failure
    2. To provide isolation between programs that access the database concurrently

All db transaction should satisfy Acid Property

ACID Property

    1. Atomicity - Either all operation complete successfully or the transaction fails 
                    and the db is unchanged.

    2. Consistency - The db state must be valid after the transacction.
                     All constraints must be satisfied

    3. Isolation - Concurrent transaction must not affect each           other 

    4. Durabiligy - Data written by a successful transaction must be recorded in persistent storage                               

How ro run SQL TX
success
BEGIN;
...
COMMIT;

error 
BEGIN;
...
ROLLBACK;