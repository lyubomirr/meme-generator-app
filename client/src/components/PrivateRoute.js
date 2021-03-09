import React from 'react';
import { Route, Redirect } from 'react-router-dom';
import { getCurrentUser } from '../auth';

const PrivateRoute = ({component: Component, user, ...rest}) => {
    return (
        <Route {...rest} render={props => (
            getCurrentUser() ?
                <Component {...props} user={user} />
            : <Redirect to="/login" />
        )} />
    );
};

export default PrivateRoute;