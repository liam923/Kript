import React, { Component } from 'react';

import Aux from '../../hoc/Aux/Aux';

class Layout extends Component {
    render() {
        return (
            <Aux>
                <div style={{height: '100px'}}></div>
                <main>
                    {this.props.children}
                </main>
            </Aux>
        );
    }
}

export default Layout;