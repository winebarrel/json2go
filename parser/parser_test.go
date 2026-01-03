package parser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/winebarrel/json2go/v2/parser"
)

func TestParseJSON_ParseErr(t *testing.T) {
	_, err := parser.ParseJSON("", []byte(`{`))
	assert.ErrorContains(t, err, `1:2: unexpected token "<EOF>" (expected "}")`)
}

func TestParseJSON_LexErr(t *testing.T) {
	_, err := parser.ParseJSON("", []byte(`{"foo:"bar"}`))
	assert.ErrorContains(t, err, `1:8: invalid character 'b' after object key`)
}

func TestParseJSON_OK(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected *parser.JsonValue
	}{
		{
			name: "int",
			json: "1",
			expected: &parser.JsonValue{
				Number: ptr("1"),
			},
		},
		{
			name: "float",
			json: "1.1",
			expected: &parser.JsonValue{
				Number: ptr("1.1"),
			},
		},
		{
			name: "false",
			json: "false",
			expected: &parser.JsonValue{
				False: ptr("false"),
			},
		},
		{
			name: "null",
			json: "null",
			expected: &parser.JsonValue{
				Null: ptr("null"),
			},
		},
		{
			name: "true",
			json: "true",
			expected: &parser.JsonValue{
				True: ptr("true"),
			},
		},
		{
			name: "string",
			json: `"hello"`,
			expected: &parser.JsonValue{
				String: ptr("hello"),
			},
		},
		{
			name: "true-string",
			json: `"true"`,
			expected: &parser.JsonValue{
				String: ptr("true"),
			},
		},
		{
			name: "false-string",
			json: `"false"`,
			expected: &parser.JsonValue{
				String: ptr("false"),
			},
		},
		{
			name: "empty object",
			json: "{}",
			expected: &parser.JsonValue{
				Object: &parser.JsonObject{},
			},
		},
		{
			name: "object",
			json: `{"str":"s","num":1,"t":true,"f":false,"null":null,"obj":{"str":"s","num":1,"t":true,"f":false,"null":null},"ary":["s",1,true,false,null]}`,
			expected: &parser.JsonValue{
				Object: &parser.JsonObject{
					Members: []*parser.JsonObjectMember{
						{
							Key: "str",
							Value: &parser.JsonValue{
								String: ptr("s"),
							},
						},
						{
							Key: "num",
							Value: &parser.JsonValue{
								Number: ptr("1"),
							},
						},
						{
							Key: "t",
							Value: &parser.JsonValue{
								True: ptr("true"),
							},
						},
						{
							Key: "f",
							Value: &parser.JsonValue{
								False: ptr("false"),
							},
						},
						{
							Key: "null",
							Value: &parser.JsonValue{
								Null: ptr("null"),
							},
						},
						{
							Key: "obj",
							Value: &parser.JsonValue{
								Object: &parser.JsonObject{
									Members: []*parser.JsonObjectMember{
										{
											Key: "str",
											Value: &parser.JsonValue{
												String: ptr("s"),
											},
										},
										{
											Key: "num",
											Value: &parser.JsonValue{
												Number: ptr("1"),
											},
										},
										{
											Key: "t",
											Value: &parser.JsonValue{
												True: ptr("true"),
											},
										},
										{
											Key: "f",
											Value: &parser.JsonValue{
												False: ptr("false"),
											},
										},
										{
											Key: "null",
											Value: &parser.JsonValue{
												Null: ptr("null"),
											},
										},
									},
								},
							},
						},
						{
							Key: "ary",
							Value: &parser.JsonValue{
								Array: &parser.JsonArray{
									Elements: []*parser.JsonValue{
										{
											String: ptr("s"),
										},
										{
											Number: ptr("1"),
										},
										{
											True: ptr("true"),
										},
										{
											False: ptr("false"),
										},
										{
											Null: ptr("null"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "array in object",
			json: `{"str":"s","num":1,"t":true,"f":false,"null":null,"ary":[{"str":"s"},{"num":1},{"t":true},{"f":false},{"null":null},[{"str":"s"},{"num":1},{"t":true},{"f":false},{"null":null}]]}`,
			expected: &parser.JsonValue{
				Object: &parser.JsonObject{
					Members: []*parser.JsonObjectMember{
						{
							Key: "str",
							Value: &parser.JsonValue{
								String: ptr("s"),
							},
						},
						{
							Key: "num",
							Value: &parser.JsonValue{
								Number: ptr("1"),
							},
						},
						{
							Key: "t",
							Value: &parser.JsonValue{
								True: ptr("true"),
							},
						},
						{
							Key: "f",
							Value: &parser.JsonValue{
								False: ptr("false"),
							},
						},
						{
							Key: "null",
							Value: &parser.JsonValue{
								Null: ptr("null"),
							},
						},
						{
							Key: "ary",
							Value: &parser.JsonValue{
								Array: &parser.JsonArray{
									Elements: []*parser.JsonValue{
										{
											Object: &parser.JsonObject{
												Members: []*parser.JsonObjectMember{
													{
														Key: "str",
														Value: &parser.JsonValue{
															String: ptr("s"),
														},
													},
												},
											},
										},
										{
											Object: &parser.JsonObject{
												Members: []*parser.JsonObjectMember{
													{
														Key: "num",
														Value: &parser.JsonValue{
															Number: ptr("1"),
														},
													},
												},
											},
										},
										{
											Object: &parser.JsonObject{
												Members: []*parser.JsonObjectMember{
													{
														Key: "t",
														Value: &parser.JsonValue{
															True: ptr("true"),
														},
													},
												},
											},
										},
										{
											Object: &parser.JsonObject{
												Members: []*parser.JsonObjectMember{
													{
														Key: "f",
														Value: &parser.JsonValue{
															False: ptr("false"),
														},
													},
												},
											},
										},
										{
											Object: &parser.JsonObject{
												Members: []*parser.JsonObjectMember{
													{
														Key: "null",
														Value: &parser.JsonValue{
															Null: ptr("null"),
														},
													},
												},
											},
										},
										{
											Array: &parser.JsonArray{
												Elements: []*parser.JsonValue{
													{
														Object: &parser.JsonObject{
															Members: []*parser.JsonObjectMember{
																{
																	Key: "str",
																	Value: &parser.JsonValue{
																		String: ptr("s"),
																	},
																},
															},
														},
													},
													{
														Object: &parser.JsonObject{
															Members: []*parser.JsonObjectMember{
																{
																	Key: "num",
																	Value: &parser.JsonValue{
																		Number: ptr("1"),
																	},
																},
															},
														},
													},
													{
														Object: &parser.JsonObject{
															Members: []*parser.JsonObjectMember{
																{
																	Key: "t",
																	Value: &parser.JsonValue{
																		True: ptr("true"),
																	},
																},
															},
														},
													},
													{
														Object: &parser.JsonObject{
															Members: []*parser.JsonObjectMember{
																{
																	Key: "f",
																	Value: &parser.JsonValue{
																		False: ptr("false"),
																	},
																},
															},
														},
													},
													{
														Object: &parser.JsonObject{
															Members: []*parser.JsonObjectMember{
																{
																	Key: "null",
																	Value: &parser.JsonValue{
																		Null: ptr("null"),
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "empty array",
			json: "[]",
			expected: &parser.JsonValue{
				Array: &parser.JsonArray{},
			},
		},
		{
			name: "object in array",
			json: `[{"str":"s"},{"num":1},{"t":true},{"f":false},{"null":null},[{"str":"s"},{"num":1},{"t":true},{"f":false},{"null":null}]]`,
			expected: &parser.JsonValue{
				Array: &parser.JsonArray{
					Elements: []*parser.JsonValue{
						{
							Object: &parser.JsonObject{
								Members: []*parser.JsonObjectMember{
									{
										Key: "str",
										Value: &parser.JsonValue{
											String: ptr("s"),
										},
									},
								},
							},
						},
						{
							Object: &parser.JsonObject{
								Members: []*parser.JsonObjectMember{
									{
										Key: "num",
										Value: &parser.JsonValue{
											Number: ptr("1"),
										},
									},
								},
							},
						},
						{
							Object: &parser.JsonObject{
								Members: []*parser.JsonObjectMember{
									{
										Key: "t",
										Value: &parser.JsonValue{
											True: ptr("true"),
										},
									},
								},
							},
						},
						{
							Object: &parser.JsonObject{
								Members: []*parser.JsonObjectMember{
									{
										Key: "f",
										Value: &parser.JsonValue{
											False: ptr("false"),
										},
									},
								},
							},
						},
						{
							Object: &parser.JsonObject{
								Members: []*parser.JsonObjectMember{
									{
										Key: "null",
										Value: &parser.JsonValue{
											Null: ptr("null"),
										},
									},
								},
							},
						},
						{
							Array: &parser.JsonArray{
								Elements: []*parser.JsonValue{
									{
										Object: &parser.JsonObject{
											Members: []*parser.JsonObjectMember{
												{
													Key: "str",
													Value: &parser.JsonValue{
														String: ptr("s"),
													},
												},
											},
										},
									},
									{
										Object: &parser.JsonObject{
											Members: []*parser.JsonObjectMember{
												{
													Key: "num",
													Value: &parser.JsonValue{
														Number: ptr("1"),
													},
												},
											},
										},
									},
									{
										Object: &parser.JsonObject{
											Members: []*parser.JsonObjectMember{
												{
													Key: "t",
													Value: &parser.JsonValue{
														True: ptr("true"),
													},
												},
											},
										},
									},
									{
										Object: &parser.JsonObject{
											Members: []*parser.JsonObjectMember{
												{
													Key: "f",
													Value: &parser.JsonValue{
														False: ptr("false"),
													},
												},
											},
										},
									},
									{
										Object: &parser.JsonObject{
											Members: []*parser.JsonObjectMember{
												{
													Key: "null",
													Value: &parser.JsonValue{
														Null: ptr("null"),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := parser.ParseJSON("", []byte(tt.json))
			require.NoError(t, err)
			assert.Equal(t, tt.expected, v)
		})
	}
}
