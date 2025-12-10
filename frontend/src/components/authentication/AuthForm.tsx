import React from 'react';

// Props for the AuthForm component
interface AuthFormProps {
    heading: string,
    buttonLabel: string,
    onSubmit: (username: string) => Promise<void>;
    extraLink: React.ReactNode;
}

// Reusable authentication form component that is used for both login and signup since both only require username
export default function AuthForm({ heading, buttonLabel, onSubmit, extraLink }: AuthFormProps) {

    // State, ensures proper validation messages are shown
    const [nameError, setNameError] = React.useState(false);
    const [nameErrorMessage, setNameErrorMessage] = React.useState('');

    // validates inputs and sets error messages if needed
    const validateInputs = () => {
        const name = document.getElementById('name') as HTMLInputElement;
        let isValid = true;

        if (!name.value || name.value.length < 1) {
            setNameError(true);
            setNameErrorMessage('Username is required.');
            isValid = false;
        } else {
            setNameError(false);
            setNameErrorMessage('');
        }

        return isValid;
    };

    // handles form submission 
    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        if (!validateInputs()) {
            return;
        }
        const username = new FormData(event.currentTarget).get('name') as string;
        await onSubmit(username);
    };

    return (
        <div>
            <h1>{heading}</h1>

            <form onSubmit={handleSubmit}>
                <label htmlFor="name">Username</label>
                <input type="text" id="name" name="name" />
                {nameError && <p>{nameErrorMessage}</p>}
                <button type="submit">{buttonLabel}</button>
                {extraLink}
            </form>
        </div>
    );
}
