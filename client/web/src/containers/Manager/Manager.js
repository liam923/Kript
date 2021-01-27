import React, { Component } from 'react';

import Aux from '../../hoc/Aux/Aux';
import Datum from '../../components/Datum/Datum';
import NavigationBar from '../../components/NavigationBar/NavigationBar';

import styles from './Manager.module.css';

class Manager extends Component {
    state = {
        info: {
            passwords: [
                { url: 'www.amazon.com', username: 'abc@gmail.com', password: 'password123' },
                { url: 'www.fb.com', username: 'abc@gmail.com', password: 'password123' },
                { url: 'www.gmail.com', username: 'abc@gmail.com', password: 'jaimeisfoxy' },
            ],
            creditcards: [
                { number: '4111 1111 1111 1111', name: 'Andy B. Calvin', expMonth: '02', expYear: '21', cvv: '123', desc: 'Visa' },
                { number: '4111 1111 1111 4444', name: 'Andy B. Calvin', expMonth: '04', expYear: '25', cvv: '123', desc: 'Mastercard' },
            ],
            codes: [
                { code: '1234', type: 'Bank Account', desc: 'Bank Account' }
            ],
            notes: [],
        },
        filter: '',
        selectedDatumId: null,
        adding: false,
    }

    // is there a better way of handling this using React internal tools?
    filterListHandler = (key) => {
        this.setState({ filter: key });
    }

    showDataHandler = (datum) => {
        console.log(datum);
        this.setState({ selectedDatumId: datum });
    }

    addDataHandler = () => {
        this.setState({ adding: true });
    }

    logOutHandler = () => {
        alert('Logged out!');
    }

    render() {
        let items = Object.keys(this.state.info).filter((pKey) => (this.state.filter === '' ? true : pKey === this.state.filter))
            .map(pKey => {
                return this.state.info[pKey].map((el, i) => {
                    return <Datum key={pKey + i} type={pKey} datum={el} clicked={this.showDataHandler} />
                });
            }).reduce((arr, el) => {
                return arr.concat(el);
            }, []);

        return (
            <Aux>
                <div style={{ display: 'flex' }}>
                    <div style={{ width: '200px' }}>
                        <div id='spacer' style={{ width: '100%', height: '118px' }}></div>
                        <NavigationBar
                            filter={this.state.filter}
                            itemClicked={this.filterListHandler}
                            addDataClicked={this.addDataHandler}
                            logOutClicked={this.logOutHandler} />
                    </div>
                    <div style={{ marginLeft: '32px', flexGrow: '1' }}>
                        <h1>Kript</h1>
                        <h2>My Data</h2>
                        {items.length !== 0 ?
                            <div className={styles.Manager}>
                                {items}
                            </div> : 
                            <p style={{ textAlign: 'center' }}>Nothing to see here.</p>}
                    </div>
                </div>
            </Aux>
        );
    }
}

export default Manager;